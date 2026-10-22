package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"multicatalogo-backend/domain"
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
)

// CheckoutRepository contiene los puertos de persistencia y caché del checkout.
type CheckoutRepository struct {
	db    interfaces.DBRepository
	cache interfaces.CacheRepository
}

func NewCheckoutRepository(db interfaces.DBRepository, cache interfaces.CacheRepository) *CheckoutRepository {
	return &CheckoutRepository{db: db, cache: cache}
}

// CrearPedido recalcula precios desde PostgreSQL y guarda pedido, stock y ledger en una sola transacción.
func (r *CheckoutRepository) CrearPedido(ctx context.Context, email, metodoPago string, items []models.CarritoItem) (*models.Pedido, error) {
	if r.db == nil {
		return nil, errors.New("base de datos no disponible")
	}
	if len(items) == 0 {
		return nil, errors.New("el carrito no puede estar vacío")
	}
	tx, err := r.db.Begin()
	if err != nil {
		return nil, fmt.Errorf("iniciar transacción: %w", err)
	}
	committed := false
	defer func() {
		if !committed {
			_ = tx.Rollback()
		}
	}()

	var usuarioID int
	if err = tx.QueryRow(`SELECT id FROM usuarios WHERE email = $1 AND activo = true`, email).Scan(&usuarioID); err != nil {
		return nil, fmt.Errorf("usuario no válido para checkout: %w", err)
	}

	type pricedItem struct {
		id, cantidad int
		precio       float64
	}
	priced := make([]pricedItem, 0, len(items))
	total := 0.0
	for _, item := range items {
		if item.ProductoID <= 0 || item.Cantidad <= 0 {
			return nil, errors.New("producto y cantidad deben ser válidos")
		}
		var precio float64
		var stock int
		err = tx.QueryRow(`SELECT precio, stock FROM productos WHERE id = $1 AND activo = true FOR UPDATE`, item.ProductoID).Scan(&precio, &stock)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("producto %d no encontrado", item.ProductoID)
			}
			return nil, fmt.Errorf("consultar producto %d: %w", item.ProductoID, err)
		}
		if stock < item.Cantidad {
			return nil, fmt.Errorf("stock insuficiente para producto %d", item.ProductoID)
		}
		priced = append(priced, pricedItem{id: item.ProductoID, cantidad: item.Cantidad, precio: precio})
		total += precio * float64(item.Cantidad)
	}

	fecha := time.Now()
	var pedidoID int
	err = tx.QueryRow(`INSERT INTO pedidos (email, usuario_id, total, metodo_pago, estado, fecha_creacion)
		VALUES ($1, $2, $3, $4, 'pendiente', $5) RETURNING id`, email, usuarioID, total, metodoPago, fecha).Scan(&pedidoID)
	if err != nil {
		return nil, fmt.Errorf("insertar pedido: %w", err)
	}
	for _, item := range priced {
		if err = tx.Exec(`INSERT INTO pedido_items (pedido_id, producto_id, cantidad, precio_unitario, subtotal)
			VALUES ($1, $2, $3, $4, $5)`, pedidoID, item.id, item.cantidad, item.precio, item.precio*float64(item.cantidad)); err != nil {
			return nil, fmt.Errorf("insertar detalle: %w", err)
		}
		if err = tx.Exec(`UPDATE productos SET stock = stock - $1 WHERE id = $2 AND stock >= $1`, item.cantidad, item.id); err != nil {
			return nil, fmt.Errorf("descontar stock: %w", err)
		}
	}
	if err = r.insertarComisiones(tx, pedidoID, email, total); err != nil {
		return nil, fmt.Errorf("registrar comisiones: %w", err)
	}
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("confirmar pedido: %w", err)
	}
	committed = true

	return &models.Pedido{ID: pedidoID, Email: email, Total: total, MetodoPago: metodoPago, Estado: "pendiente", FechaCreacion: fecha.Format(time.RFC3339)}, nil
}

// insertarComisiones escribe el ledger con las tasas 10%, 5% y 2% antes del COMMIT.
func (r *CheckoutRepository) insertarComisiones(tx interfaces.Transaction, pedidoID int, email string, total float64) error {
	var referidoID sql.NullInt64
	if err := tx.QueryRow(`SELECT parent_id FROM referidos WHERE email = $1`, email).Scan(&referidoID); err != nil {
		if errors.Is(err, sql.ErrNoRows) { // Cuenta válida sin nodo de red.
			return nil
		}
		return err
	}
	porcentajes := []float64{10, 5, 2}
	comisiones := domain.RepartirComisiones(total, porcentajes)
	current := referidoID
	for _, comision := range comisiones {
		if !current.Valid {
			break
		}
		id := int(current.Int64)
		var parent sql.NullInt64
		if err := tx.QueryRow(`SELECT parent_id FROM referidos WHERE id = $1`, id).Scan(&parent); err != nil {
			return err
		}
		if err := tx.Exec(`INSERT INTO comisiones (pedido_id, referido_id, nivel, porcentaje, monto, estado)
			VALUES ($1, $2, $3, $4, $5, 'pendiente')`, pedidoID, id, comision.Nivel, comision.Porcentaje, comision.Monto); err != nil {
			return err
		}
		current = parent
	}
	return nil
}

// CalcularComisiones actualiza los contadores después del COMMIT; el ledger ya fue persistido atómicamente.
func (r *CheckoutRepository) CalcularComisiones(ctx context.Context, email string, total float64) error {
	if r.db == nil {
		return errors.New("base de datos no disponible")
	}
	row, err := r.db.QueryRow(`SELECT parent_id FROM referidos WHERE email = $1`, email)
	if err != nil {
		return err
	}
	current := nullableInt(row["parent_id"])
	for i := 0; i < 3 && current.Valid; i++ {
		id := int(current.Int64)
		ancestor, err := r.db.QueryRow(`SELECT parent_id FROM referidos WHERE id = $1`, id)
		if err != nil {
			return err
		}
		if err = r.db.Exec(`UPDATE referidos SET ventas = COALESCE(ventas, 0) + $1 WHERE id = $2`, total, id); err != nil {
			return err
		}
		current = nullableInt(ancestor["parent_id"])
	}
	if r.cache != nil {
		_ = r.cache.Delete(ctx, "productos:all")
		_ = r.cache.Delete(ctx, fmt.Sprintf("red:tree:%s", email))
	}
	return nil
}

func nullableInt(value interface{}) sql.NullInt64 {
	switch v := value.(type) {
	case int64:
		return sql.NullInt64{Int64: v, Valid: true}
	case int:
		return sql.NullInt64{Int64: int64(v), Valid: true}
	case []byte:
		var n int64
		if _, err := fmt.Sscan(string(v), &n); err == nil {
			return sql.NullInt64{Int64: n, Valid: true}
		}
	}
	return sql.NullInt64{}
}
