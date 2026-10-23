package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
)

type BackofficeController struct {
	db    *sql.DB
	cache interfaces.CacheRepository
}

func NewBackofficeController(db *sql.DB, cache interfaces.CacheRepository) *BackofficeController {
	return &BackofficeController{db: db, cache: cache}
}

type productInput struct {
	Nombre      string   `json:"nombre"`
	Descripcion string   `json:"descripcion"`
	Precio      float64  `json:"precio"`
	Categoria   string   `json:"categoria"`
	Img         string   `json:"img"`
	Galeria     []string `json:"galeria"`
	Stock       int      `json:"stock"`
}

func pagination(c *fiber.Ctx) (int, int, int) {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	return page, limit, (page - 1) * limit
}

func (bc *BackofficeController) CreateProduct(c *fiber.Ctx) error {
	var in productInput
	if err := c.BodyParser(&in); err != nil {
		return apiError(c, 400, "Cuerpo inválido", nil)
	}
	if strings.TrimSpace(in.Nombre) == "" || in.Precio <= 0 || strings.TrimSpace(in.Categoria) == "" || in.Stock < 0 {
		return apiError(c, 400, "nombre, precio, categoría y stock son inválidos", nil)
	}
	gallery, _ := json.Marshal(in.Galeria)
	var p models.Producto
	err := bc.db.QueryRow(`INSERT INTO productos(nombre,descripcion,precio,categoria,img,galeria,stock) VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`, in.Nombre, in.Descripcion, in.Precio, in.Categoria, in.Img, gallery, in.Stock).Scan(&p.ID)
	if err != nil {
		return apiError(c, 409, "No se pudo crear el producto", nil)
	}
	p.Nombre, p.Descripcion, p.Precio, p.Categoria, p.Img, p.Galeria, p.Stock = in.Nombre, in.Descripcion, in.Precio, in.Categoria, in.Img, in.Galeria, in.Stock
	bc.invalidateProducts()
	return c.Status(201).JSON(p)
}

func (bc *BackofficeController) UpdateProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apiError(c, 400, "ID inválido", nil)
	}
	var in productInput
	if err = c.BodyParser(&in); err != nil {
		return apiError(c, 400, "Cuerpo inválido", nil)
	}
	if in.Precio < 0 || in.Stock < 0 {
		return apiError(c, 400, "precio o stock inválido", nil)
	}
	gallery, _ := json.Marshal(in.Galeria)
	res, err := bc.db.Exec(`UPDATE productos SET nombre=COALESCE(NULLIF($1,''),nombre), descripcion=COALESCE(NULLIF($2,''),descripcion), precio=CASE WHEN $3>0 THEN $3 ELSE precio END, categoria=COALESCE(NULLIF($4,''),categoria), img=COALESCE(NULLIF($5,''),img), galeria=CASE WHEN $6::jsonb='null'::jsonb THEN galeria ELSE $6::jsonb END, stock=$7 WHERE id=$8`, in.Nombre, in.Descripcion, in.Precio, in.Categoria, in.Img, gallery, in.Stock, id)
	if err != nil {
		return apiError(c, 500, "No se pudo actualizar el producto", nil)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apiError(c, 404, "Producto no encontrado", nil)
	}
	bc.invalidateProducts()
	return c.JSON(fiber.Map{"message": "Producto actualizado"})
}

func (bc *BackofficeController) DeleteProduct(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apiError(c, 400, "ID inválido", nil)
	}
	res, err := bc.db.Exec(`UPDATE productos SET activo=false WHERE id=$1`, id)
	if err != nil {
		return apiError(c, 500, "No se pudo desactivar", nil)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return apiError(c, 404, "Producto no encontrado", nil)
	}
	bc.invalidateProducts()
	return c.SendStatus(204)
}

func (bc *BackofficeController) AdminProducts(c *fiber.Ctx) error {
	page, limit, offset := pagination(c)
	var total int
	if err := bc.db.QueryRow(`SELECT COUNT(*) FROM productos`).Scan(&total); err != nil {
		return apiError(c, 500, "Error al contar productos", nil)
	}
	rows, err := bc.db.Query(`SELECT id,nombre,descripcion,precio,categoria,img,galeria,stock,activo FROM productos ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return apiError(c, 500, "Error al consultar productos", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id, stock int
		var name, desc, cat, img string
		var price float64
		var gallery []byte
		var active bool
		if err = rows.Scan(&id, &name, &desc, &price, &cat, &img, &gallery, &stock, &active); err != nil {
			return apiError(c, 500, "Error al leer productos", nil)
		}
		var g []string
		_ = json.Unmarshal(gallery, &g)
		out = append(out, fiber.Map{"id": id, "nombre": name, "descripcion": desc, "precio": price, "categoria": cat, "img": img, "galeria": g, "stock": stock, "activo": active})
	}
	return c.JSON(fiber.Map{"data": out, "page": page, "limit": limit, "total": total})
}

func (bc *BackofficeController) AdminUsers(c *fiber.Ctx) error {
	page, limit, offset := pagination(c)
	var total int
	if err := bc.db.QueryRow(`SELECT COUNT(*) FROM usuarios`).Scan(&total); err != nil {
		return apiError(c, 500, "Error al contar usuarios", nil)
	}
	rows, err := bc.db.Query(`SELECT id,email,nombre,rol,activo,creado_en FROM usuarios ORDER BY id LIMIT $1 OFFSET $2`, limit, offset)
	if err != nil {
		return apiError(c, 500, "Error al consultar usuarios", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id int
		var email, name, role string
		var active bool
		var created time.Time
		if err := rows.Scan(&id, &email, &name, &role, &active, &created); err != nil {
			return apiError(c, 500, "Error al leer usuarios", nil)
		}
		out = append(out, fiber.Map{"id": id, "email": email, "nombre": name, "rol": role, "activo": active, "creado_en": created})
	}
	return c.JSON(fiber.Map{"data": out, "page": page, "limit": limit, "total": total})
}

func (bc *BackofficeController) UpdateUserState(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apiError(c, 400, "ID inválido", nil)
	}
	var in struct {
		Activo *bool `json:"activo"`
	}
	if err = c.BodyParser(&in); err != nil || in.Activo == nil {
		return apiError(c, 400, "El estado activo es obligatorio", nil)
	}
	if id == 1 && !*in.Activo {
		return apiError(c, 400, "No se puede desactivar al administrador principal", nil)
	}
	res, err := bc.db.Exec(`UPDATE usuarios SET activo=$1 WHERE id=$2`, *in.Activo, id)
	if err != nil {
		return apiError(c, 500, "No se pudo actualizar el usuario", nil)
	}
	count, _ := res.RowsAffected()
	if count == 0 {
		return apiError(c, 404, "Usuario no encontrado", nil)
	}
	return c.JSON(fiber.Map{"message": "Estado actualizado", "activo": *in.Activo})
}

func (bc *BackofficeController) MyOrders(c *fiber.Ctx) error {
	email, _ := c.Locals("email").(string)
	return bc.listOrders(c, `WHERE p.email=$1`, email)
}
func (bc *BackofficeController) AdminOrders(c *fiber.Ctx) error {
	if state := c.Query("estado"); state != "" {
		return bc.listOrders(c, `WHERE p.estado=$1`, state)
	}
	return bc.listOrders(c, "", nil)
}
func (bc *BackofficeController) listOrders(c *fiber.Ctx, where string, arg interface{}) error {
	page, limit, offset := pagination(c)
	countQuery := `SELECT COUNT(*) FROM pedidos p`
	args := make([]interface{}, 0, 1)
	if where != "" {
		countQuery += " " + where
		args = append(args, arg)
	}
	var total int
	if err := bc.db.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return apiError(c, 500, "Error al contar pedidos", nil)
	}
	q := `SELECT p.id,p.email,p.total,p.metodo_pago,p.estado,p.fecha_creacion,COUNT(i.id) FROM pedidos p LEFT JOIN pedido_items i ON i.pedido_id=p.id`
	if where != "" {
		q += " " + where
	}
	q += fmt.Sprintf(` GROUP BY p.id ORDER BY p.fecha_creacion DESC LIMIT $%d OFFSET $%d`, len(args)+1, len(args)+2)
	args = append(args, limit, offset)
	var rows *sql.Rows
	var err error
	rows, err = bc.db.Query(q, args...)
	if err != nil {
		return apiError(c, 500, "Error al consultar pedidos", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id, count int
		var email, method, state string
		var total float64
		var date time.Time
		if err = rows.Scan(&id, &email, &total, &method, &state, &date, &count); err != nil {
			return apiError(c, 500, "Error al leer pedidos", nil)
		}
		out = append(out, fiber.Map{"id": id, "email": email, "total": total, "metodo_pago": method, "estado": state, "fecha_creacion": date, "items": count})
	}
	return c.JSON(fiber.Map{"data": out, "page": page, "limit": limit, "total": total})
}

func (bc *BackofficeController) OrderDetail(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apiError(c, 400, "ID inválido", nil)
	}
	email, _ := c.Locals("email").(string)
	role, _ := c.Locals("rol").(string)
	var owner, method, state string
	var total float64
	var date time.Time
	err = bc.db.QueryRow(`SELECT id,email,total,metodo_pago,estado,fecha_creacion FROM pedidos WHERE id=$1`, id).Scan(&id, &owner, &total, &method, &state, &date)
	if err != nil {
		return apiError(c, 404, "Pedido no encontrado", nil)
	}
	if role != "admin" && owner != email {
		return apiError(c, 403, "No puedes consultar este pedido", nil)
	}
	var items []fiber.Map
	rows, err := bc.db.Query(`SELECT i.producto_id,i.cantidad,i.precio_unitario,i.subtotal,COALESCE(p.nombre,'')
FROM pedido_items i LEFT JOIN productos p ON p.id=i.producto_id WHERE i.pedido_id=$1`, id)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var pid, q int
			var price, subtotal float64
			var name string
			_ = rows.Scan(&pid, &q, &price, &subtotal, &name)
			items = append(items, fiber.Map{"producto_id": pid, "nombre": name, "cantidad": q, "precio_unitario": price, "subtotal": subtotal})
		}
	}
	return c.JSON(fiber.Map{"id": id, "email": owner, "total": total, "metodo_pago": method, "estado": state, "fecha_creacion": date, "items": items})
}

func (bc *BackofficeController) UpdateOrderState(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return apiError(c, 400, "ID inválido", nil)
	}
	var in struct {
		Estado string `json:"estado"`
	}
	if err = c.BodyParser(&in); err != nil {
		return apiError(c, 400, "Cuerpo inválido", nil)
	}
	allowed := map[string]bool{"pendiente": true, "preparando": true, "enviado": true, "entregado": true, "cancelado": true}
	if !allowed[in.Estado] {
		return apiError(c, 400, "Estado inválido", nil)
	}
	tx, err := bc.db.Begin()
	if err != nil {
		return apiError(c, 500, "No se pudo iniciar la operación", nil)
	}
	defer tx.Rollback()
	var old string
	if err = tx.QueryRow(`SELECT estado FROM pedidos WHERE id=$1 FOR UPDATE`, id).Scan(&old); err != nil {
		return apiError(c, 404, "Pedido no encontrado", nil)
	}
	valid := old == "pendiente" && (in.Estado == "preparando" || in.Estado == "cancelado") || old == "preparando" && in.Estado == "enviado" || old == "enviado" && in.Estado == "entregado"
	if !valid {
		return apiError(c, 400, "Transición de estado inválida", nil)
	}
	if in.Estado == "cancelado" {
		rows, _ := tx.Query(`SELECT producto_id,cantidad FROM pedido_items WHERE pedido_id=$1`, id)
		for rows.Next() {
			var pid, q int
			_ = rows.Scan(&pid, &q)
			_, err = tx.Exec(`UPDATE productos SET stock=stock+$1 WHERE id=$2`, q, pid)
			if err != nil {
				return apiError(c, 500, "No se pudo restaurar stock", nil)
			}
		}
		_ = rows.Close()
		if _, err = tx.Exec(`UPDATE comisiones SET estado='revertida' WHERE pedido_id=$1`, id); err != nil {
			return apiError(c, 500, "No se pudieron revertir comisiones", nil)
		}
	}
	if _, err = tx.Exec(`UPDATE pedidos SET estado=$1 WHERE id=$2`, in.Estado, id); err != nil {
		return apiError(c, 500, "No se pudo actualizar pedido", nil)
	}
	if err = tx.Commit(); err != nil {
		return apiError(c, 500, "No se pudo confirmar", nil)
	}
	bc.invalidateProducts()
	return c.JSON(fiber.Map{"message": "Estado actualizado", "estado": in.Estado})
}

func (bc *BackofficeController) MyCommissions(c *fiber.Ctx) error {
	email, _ := c.Locals("email").(string)
	rows, err := bc.db.Query(`SELECT c.id,c.pedido_id,c.nivel,c.porcentaje,c.monto,c.estado,c.creado_en FROM comisiones c JOIN referidos r ON r.id=c.referido_id WHERE r.email=$1 ORDER BY c.creado_en DESC`, email)
	if err != nil {
		return apiError(c, 500, "Error al consultar comisiones", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	var pending, paid float64
	for rows.Next() {
		var id, pid, level int
		var pct, monto float64
		var state string
		var created time.Time
		if err = rows.Scan(&id, &pid, &level, &pct, &monto, &state, &created); err != nil {
			return apiError(c, 500, "Error al leer comisiones", nil)
		}
		if state == "pendiente" {
			pending += monto
		}
		if state == "pagada" {
			paid += monto
		}
		out = append(out, fiber.Map{"id": id, "pedido_id": pid, "nivel": level, "porcentaje": pct, "monto": monto, "estado": state, "creado_en": created})
	}
	return c.JSON(fiber.Map{"data": out, "pendiente": pending, "pagada": paid})
}

func (bc *BackofficeController) Liquidate(c *fiber.Ctx) error {
	var in struct {
		ReferidoID int `json:"referido_id"`
	}
	if err := c.BodyParser(&in); err != nil || in.ReferidoID <= 0 {
		return apiError(c, 400, "referido_id inválido", nil)
	}
	tx, err := bc.db.Begin()
	if err != nil {
		return apiError(c, 500, "No se pudo iniciar liquidación", nil)
	}
	defer tx.Rollback()
	rows, err := tx.Query(`SELECT monto FROM comisiones WHERE referido_id=$1 AND estado='pendiente' FOR UPDATE`, in.ReferidoID)
	if err != nil {
		return apiError(c, 500, "Error al calcular deuda", nil)
	}
	var total float64
	for rows.Next() {
		var amount float64
		if err := rows.Scan(&amount); err != nil {
			rows.Close()
			return apiError(c, 500, "Error al leer deuda", nil)
		}
		total += amount
	}
	rows.Close()
	if total <= 0 {
		return apiError(c, 409, "No hay comisiones pendientes", nil)
	}
	var id int
	if err = tx.QueryRow(`INSERT INTO liquidaciones(referido_id,total,estado,pagada_en) VALUES($1,$2,'pagada',NOW()) RETURNING id`, in.ReferidoID, total).Scan(&id); err != nil {
		return apiError(c, 500, "Error al crear liquidación", nil)
	}
	if _, err = tx.Exec(`UPDATE comisiones SET estado='pagada',liquidacion_id=$1 WHERE referido_id=$2 AND estado='pendiente'`, id, in.ReferidoID); err != nil {
		return apiError(c, 500, "Error al marcar comisiones", nil)
	}
	if err = tx.Commit(); err != nil {
		return apiError(c, 500, "No se pudo confirmar liquidación", nil)
	}
	return c.Status(201).JSON(fiber.Map{"id": id, "referido_id": in.ReferidoID, "total": total, "estado": "pagada"})
}
func (bc *BackofficeController) Liquidations(c *fiber.Ctx) error {
	rows, err := bc.db.Query(`SELECT l.id,l.referido_id,r.nombre,r.email,l.total,l.estado,l.creado_en,l.pagada_en FROM liquidaciones l JOIN referidos r ON r.id=l.referido_id ORDER BY l.creado_en DESC`)
	if err != nil {
		return apiError(c, 500, "Error al consultar liquidaciones", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id, rid int
		var name, email string
		var total float64
		var state string
		var created, payed sql.NullTime
		if err := rows.Scan(&id, &rid, &name, &email, &total, &state, &created, &payed); err != nil {
			return apiError(c, 500, "Error al leer liquidaciones", nil)
		}
		out = append(out, fiber.Map{"id": id, "referido_id": rid, "nombre": name, "email": email, "total": total, "estado": state, "creado_en": created.Time, "pagada_en": payed.Time})
	}
	return c.JSON(out)
}
func (bc *BackofficeController) MyLiquidations(c *fiber.Ctx) error {
	email, _ := c.Locals("email").(string)
	rows, err := bc.db.Query(`SELECT l.id,l.total,l.estado,l.creado_en,l.pagada_en FROM liquidaciones l JOIN referidos r ON r.id=l.referido_id WHERE r.email=$1 ORDER BY l.creado_en DESC`, email)
	if err != nil {
		return apiError(c, 500, "Error al consultar liquidaciones", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id int
		var total float64
		var state string
		var created, payed sql.NullTime
		_ = rows.Scan(&id, &total, &state, &created, &payed)
		out = append(out, fiber.Map{"id": id, "total": total, "estado": state, "creado_en": created.Time, "pagada_en": payed.Time})
	}
	return c.JSON(out)
}

func (bc *BackofficeController) Reports(c *fiber.Ctx) error {
	var total, avg, comm float64
	var count int
	where := `estado <> 'cancelado'`
	args := make([]interface{}, 0, 2)
	if desde := c.Query("desde"); desde != "" {
		args = append(args, desde)
		where += fmt.Sprintf(" AND fecha_creacion >= $%d", len(args))
	}
	if hasta := c.Query("hasta"); hasta != "" {
		args = append(args, hasta)
		where += fmt.Sprintf(" AND fecha_creacion < ($%d::date + INTERVAL '1 day')", len(args))
	}
	err := bc.db.QueryRow(`SELECT COALESCE(SUM(total),0),COUNT(*),COALESCE(AVG(total),0) FROM pedidos WHERE `+where, args...).Scan(&total, &count, &avg)
	if err != nil {
		return apiError(c, 500, "Error en reporte", nil)
	}
	commWhere := `estado <> 'revertida'`
	commArgs := make([]interface{}, 0, 2)
	if desde := c.Query("desde"); desde != "" {
		commArgs = append(commArgs, desde)
		commWhere += fmt.Sprintf(" AND creado_en >= $%d", len(commArgs))
	}
	if hasta := c.Query("hasta"); hasta != "" {
		commArgs = append(commArgs, hasta)
		commWhere += fmt.Sprintf(" AND creado_en < ($%d::date + INTERVAL '1 day')", len(commArgs))
	}
	_ = bc.db.QueryRow(`SELECT COALESCE(SUM(monto),0) FROM comisiones WHERE `+commWhere, commArgs...).Scan(&comm)
	return c.JSON(fiber.Map{"total_vendido": total, "pedidos": count, "ticket_promedio": avg, "comisiones_devengadas": comm})
}
func (bc *BackofficeController) TopProducts(c *fiber.Ctx) error {
	rows, err := bc.db.Query(`SELECT p.id,p.nombre,SUM(i.cantidad),SUM(i.subtotal) FROM pedido_items i JOIN productos p ON p.id=i.producto_id JOIN pedidos o ON o.id=i.pedido_id WHERE o.estado<>'cancelado' GROUP BY p.id,p.nombre ORDER BY SUM(i.cantidad) DESC`)
	if err != nil {
		return apiError(c, 500, "Error en reporte", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id, q int
		var name string
		var total float64
		_ = rows.Scan(&id, &name, &q, &total)
		out = append(out, fiber.Map{"producto_id": id, "nombre": name, "cantidad": q, "total": total})
	}
	return c.JSON(out)
}
func (bc *BackofficeController) PendingCommissions(c *fiber.Ctx) error {
	rows, err := bc.db.Query(`SELECT c.referido_id,r.nombre,r.email,COALESCE(SUM(c.monto),0)
FROM comisiones c JOIN referidos r ON r.id=c.referido_id
WHERE c.estado='pendiente' GROUP BY c.referido_id,r.nombre,r.email ORDER BY 4 DESC`)
	if err != nil {
		return apiError(c, 500, "Error en reporte", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id int
		var name, email string
		var total float64
		if err := rows.Scan(&id, &name, &email, &total); err != nil {
			return apiError(c, 500, "Error al leer deuda", nil)
		}
		out = append(out, fiber.Map{"referido_id": id, "nombre": name, "email": email, "total": total})
	}
	return c.JSON(out)
}

func (bc *BackofficeController) TopReferidos(c *fiber.Ctx) error {
	rows, err := bc.db.Query(`SELECT r.id,r.email,r.nombre,r.nivel,COALESCE(r.ventas,0),COALESCE(SUM(c.monto),0) FROM referidos r LEFT JOIN comisiones c ON c.referido_id=r.id AND c.estado<>'revertida' GROUP BY r.id,r.email,r.nombre,r.nivel,r.ventas ORDER BY 5 DESC`)
	if err != nil {
		return apiError(c, 500, "Error en reporte", nil)
	}
	defer rows.Close()
	out := make([]fiber.Map, 0)
	for rows.Next() {
		var id int
		var email, name string
		var level int
		var sales, commissions float64
		if err := rows.Scan(&id, &email, &name, &level, &sales, &commissions); err != nil {
			return apiError(c, 500, "Error al leer reporte", nil)
		}
		out = append(out, fiber.Map{"referido_id": id, "email": email, "nombre": name, "nivel": level, "ventas": sales, "comisiones": commissions})
	}
	return c.JSON(out)
}
func (bc *BackofficeController) invalidateProducts() {
	if bc.cache != nil {
		_ = bc.cache.Delete(context.Background(), "productos:all")
	}
}
