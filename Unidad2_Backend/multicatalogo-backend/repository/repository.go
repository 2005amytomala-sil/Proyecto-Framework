// Declaramos el paquete repository para las operaciones de acceso a datos.
package repository

import (
	// Importamos database/sql para manejar conexiones a bases de datos.
	"database/sql"
	// Importamos fmt para formatear strings.
	"encoding/json"
	"errors"
	"fmt"

	// Importamos nuestros modelos.
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
)

var ErrNotFound = errors.New("registro no encontrado")

// Repository encapsula las operaciones de acceso a datos.
type Repository struct {
	db *sql.DB
}

// NewRepository crea una nueva instancia de Repository.
func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Begin inicia una transacción.
func (r *Repository) Begin() (interfaces.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

// BuscarUsuarioPorEmail busca un usuario por email.
func (r *Repository) BuscarUsuarioPorEmail(email string) (*models.Usuario, error) {
	query := `SELECT email, password_hash, nombre, rol, activo FROM usuarios WHERE email = $1`

	var usuario models.Usuario
	err := r.db.QueryRow(query, email).Scan(&usuario.Email, &usuario.PasswordHash, &usuario.Nombre, &usuario.Rol, &usuario.Activo)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: usuario", ErrNotFound)
		}
		return nil, fmt.Errorf("error al buscar usuario: %w", err)
	}

	return &usuario, nil
}

// BuscarReferidoPorCodigo busca un referido por su código de invitación.
func (r *Repository) BuscarReferidoPorCodigo(codigo string) (int, error) {
	query := `SELECT id FROM referidos WHERE codigo = $1`

	var id int
	err := r.db.QueryRow(query, codigo).Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("%w: código de referido", ErrNotFound)
		}
		return 0, fmt.Errorf("error al buscar referido: %w", err)
	}

	return id, nil
}

// CrearUsuario crea un nuevo usuario en la base de datos.
func (r *Repository) CrearUsuario(email, passwordHash, nombre, rol string, referidoID *int) error {
	query := `
		INSERT INTO usuarios (email, password_hash, nombre, rol, referido_id)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, email, passwordHash, nombre, rol, referidoID)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

// CrearNodoReferido crea un nuevo nodo en la tabla referidos.
func (r *Repository) CrearNodoReferido(email, nombre, codigo string, nivel int, parentID *int) error {
	query := `
		INSERT INTO referidos (email, nombre, codigo, nivel, parent_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query, email, nombre, codigo, nivel, parentID)
	if err != nil {
		return fmt.Errorf("error al crear nodo referido: %w", err)
	}

	return nil
}

// CrearUsuarioTx crea un usuario dentro de una transacción.
func (r *Repository) CrearUsuarioTx(tx interfaces.Transaction, email, passwordHash, nombre, rol string) error {
	query := `
		INSERT INTO usuarios (email, password_hash, nombre, rol)
		VALUES ($1, $2, $3, $4)
	`
	err := tx.Exec(query, email, passwordHash, nombre, rol)
	if err != nil {
		return fmt.Errorf("error al crear usuario: %w", err)
	}

	return nil
}

// CrearNodoReferidoTx crea un nodo referido dentro de una transacción.
func (r *Repository) CrearNodoReferidoTx(tx interfaces.Transaction, email, nombre, codigo string, nivel int, parentID *int) (int, error) {
	query := `
		INSERT INTO referidos (email, nombre, codigo, nivel, parent_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`
	var id int
	err := tx.QueryRow(query, email, nombre, codigo, nivel, parentID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("error al crear nodo referido: %w", err)
	}
	return id, nil
}

func (r *Repository) ObtenerNivelReferido(id int) (int, error) {
	var nivel int
	err := r.db.QueryRow(`SELECT nivel FROM referidos WHERE id = $1`, id).Scan(&nivel)
	if err == sql.ErrNoRows {
		return 0, fmt.Errorf("%w: referido", ErrNotFound)
	}
	return nivel, err
}

func (r *Repository) VincularUsuarioReferidoTx(tx interfaces.Transaction, email string, nodoID int, parentID *int) error {
	return tx.Exec(`UPDATE usuarios SET referido_id = $1,
		referido_por = (SELECT email FROM referidos WHERE id = $2) WHERE email = $3`, nodoID, parentID, email)
}

// ObtenerProductos retorna todos los productos de la base de datos.
func (r *Repository) ObtenerProductos() ([]models.Producto, error) {
	return r.buscarProductos("")
}

func (r *Repository) BuscarProductos(queryText string) ([]models.Producto, error) {
	return r.buscarProductos(queryText)
}

func (r *Repository) buscarProductos(queryText string) ([]models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock FROM productos WHERE activo = true`
	args := []interface{}{}
	if queryText != "" {
		args = append(args, "%"+queryText+"%")
		query += ` AND (nombre ILIKE $1 OR descripcion ILIKE $1 OR categoria ILIKE $1)`
	}
	query += ` ORDER BY id`

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error al obtener productos: %w", err)
	}
	defer rows.Close()

	var productos []models.Producto
	for rows.Next() {
		var p models.Producto
		var descripcion, categoria, img sql.NullString
		var galeria []byte

		err := rows.Scan(&p.ID, &p.Nombre, &descripcion, &p.Precio, &categoria, &img, &galeria, &p.Stock)
		if err != nil {
			return nil, fmt.Errorf("error al escanear producto: %w", err)
		}
		p.Descripcion, p.Categoria = descripcion.String, categoria.String
		p.Img = img.String
		if len(galeria) > 0 {
			if err := json.Unmarshal(galeria, &p.Galeria); err != nil {
				return nil, fmt.Errorf("decodificar galería: %w", err)
			}
		}
		if p.Galeria == nil {
			p.Galeria = []string{}
		}

		productos = append(productos, p)
	}

	return productos, nil
}

// ObtenerProductoPorID retorna un producto específico por su ID.
func (r *Repository) ObtenerProductoPorID(id int) (*models.Producto, error) {
	query := `SELECT id, nombre, descripcion, precio, categoria, img, galeria, stock FROM productos WHERE id = $1 AND activo = true`

	var p models.Producto
	var descripcion, categoria, img sql.NullString
	var galeria []byte

	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Nombre, &descripcion, &p.Precio, &categoria, &img, &galeria, &p.Stock)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("%w: producto", ErrNotFound)
		}
		return nil, fmt.Errorf("error al obtener producto: %w", err)
	}
	p.Descripcion, p.Categoria, p.Img = descripcion.String, categoria.String, img.String
	if len(galeria) > 0 {
		if err := json.Unmarshal(galeria, &p.Galeria); err != nil {
			return nil, fmt.Errorf("decodificar galería: %w", err)
		}
	}
	if p.Galeria == nil {
		p.Galeria = []string{}
	}

	return &p, nil
}

// ObtenerRed construye el árbol de referidos multinivel desde la base de datos.
func (r *Repository) ObtenerRed(email string) (*models.NodoRed, error) {
	rows, err := r.db.Query(`SELECT r.id, r.email, r.nombre, r.codigo, r.ventas, r.nivel, r.parent_id, COALESCE(u.rol, 'cliente')
		FROM referidos r LEFT JOIN usuarios u ON u.email = r.email ORDER BY r.id`)
	if err != nil {
		return nil, fmt.Errorf("error al obtener referidos: %w", err)
	}
	defer rows.Close()
	nodes := make(map[int]*models.NodoRed)
	parents := make(map[int]sql.NullInt64)
	emailToID := make(map[string]int)
	for rows.Next() {
		var id int
		var nodeEmail, nombre, codigo, rol string
		var ventas float64
		var nivel int
		var parentID sql.NullInt64
		if err := rows.Scan(&id, &nodeEmail, &nombre, &codigo, &ventas, &nivel, &parentID, &rol); err != nil {
			return nil, fmt.Errorf("error al escanear referido: %w", err)
		}
		nodes[id] = &models.NodoRed{Usuario: models.Usuario{Email: nodeEmail, Nombre: nombre, Rol: rol}, Codigo: codigo, Ventas: ventas, Hijos: []models.NodoRed{}}
		parents[id] = parentID
		emailToID[nodeEmail] = id
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error al leer referidos: %w", err)
	}
	rootID, ok := emailToID[email]
	if !ok {
		return nil, fmt.Errorf("%w: usuario de red", ErrNotFound)
	}
	for id, node := range nodes {
		parent := parents[id]
		if parent.Valid {
			if parentNode, exists := nodes[int(parent.Int64)]; exists {
				parentNode.Hijos = append(parentNode.Hijos, *node)
			}
		}
	}
	return nodes[rootID], nil
}
