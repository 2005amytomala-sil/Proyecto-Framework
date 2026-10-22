// Declaramos el paquete interfaces para definir los contratos de inyección de dependencias.
package interfaces

import (
	// Importamos context para manejar el contexto de las operaciones.
	"context"
	"database/sql"
	"errors"
	"multicatalogo-backend/models"
)

var ErrNotFound = errors.New("registro no encontrado")

type AuthRepository interface {
	BuscarUsuarioPorEmail(string) (*models.Usuario, error)
	BuscarReferidoPorCodigo(string) (int, error)
	ObtenerNivelReferido(int) (int, error)
	Begin() (Transaction, error)
	CrearUsuarioTx(Transaction, string, string, string, string) error
	CrearNodoReferidoTx(Transaction, string, string, string, int, *int) (int, error)
	VincularUsuarioReferidoTx(Transaction, string, int, *int) error
}
type ProductRepository interface {
	ObtenerProductos() ([]models.Producto, error)
	ObtenerProductoPorID(int) (*models.Producto, error)
}
type RedRepository interface {
	ObtenerRed(string) (*models.NodoRed, error)
}
type CheckoutRepository interface {
	CrearPedido(context.Context, string, string, []models.CarritoItem) (*models.Pedido, error)
	CalcularComisiones(context.Context, string, float64) error
}

// CacheRepository define el contrato para operaciones de caché.
type CacheRepository interface {
	// Get obtiene un valor del caché por su clave.
	Get(ctx context.Context, key string) (string, error)
	// Set guarda un valor en el caché con un tiempo de expiración.
	Set(ctx context.Context, key string, value interface{}, expiration interface{}) error
	// Delete elimina una clave del caché.
	Delete(ctx context.Context, key string) error
}

// DBRepository define el contrato para operaciones de base de datos.
type DBRepository interface {
	// Begin inicia una transacción.
	Begin() (Transaction, error)
	// Exec ejecuta una consulta SQL sin retornar filas.
	Exec(query string, args ...interface{}) error
	// Query ejecuta una consulta SQL y retorna filas.
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)
	// QueryRow ejecuta una consulta SQL y retorna una sola fila.
	QueryRow(query string, args ...interface{}) (map[string]interface{}, error)
}

// Transaction define el contrato para operaciones transaccionales.
type Transaction interface {
	// Commit confirma la transacción.
	Commit() error
	// Rollback revierte la transacción.
	Rollback() error
	// Exec ejecuta una consulta SQL dentro de la transacción.
	Exec(query string, args ...interface{}) error
	// QueryRow ejecuta una consulta SQL dentro de la transacción y retorna una sola fila.
	QueryRow(query string, args ...interface{}) *sql.Row
}
