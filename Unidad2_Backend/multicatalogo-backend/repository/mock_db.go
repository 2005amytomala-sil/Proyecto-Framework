// Declaramos el paquete repository para las implementaciones concretas de los repositorios.
package repository

import (
	// Importamos nuestras interfaces.
	"database/sql"
	"multicatalogo-backend/interfaces"
)

// MockDB es una implementación mock de DBRepository para cuando PostgreSQL no está disponible.
type MockDB struct{}

// NewMockDB crea una nueva instancia de MockDB.
func NewMockDB() interfaces.DBRepository {
	return &MockDB{}
}

// Begin inicia una transacción (mock).
func (m *MockDB) Begin() (interfaces.Transaction, error) {
	return &MockTransaction{}, nil
}

// Exec ejecuta una consulta SQL (mock).
func (m *MockDB) Exec(query string, args ...interface{}) error {
	return nil
}

// Query ejecuta una consulta SQL y retorna filas (mock).
func (m *MockDB) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	return []map[string]interface{}{}, nil
}

// QueryRow ejecuta una consulta SQL y retorna una sola fila (mock).
func (m *MockDB) QueryRow(query string, args ...interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{}, nil
}

// MockTransaction es una implementación mock de Transaction.
type MockTransaction struct{}

// Commit confirma la transacción (mock).
func (m *MockTransaction) Commit() error {
	return nil
}

// Rollback revierte la transacción (mock).
func (m *MockTransaction) Rollback() error {
	return nil
}

// Exec ejecuta una consulta SQL dentro de la transacción (mock).
func (m *MockTransaction) Exec(query string, args ...interface{}) error {
	return nil
}

// QueryRow ejecuta una consulta SQL dentro de la transacción y retorna una sola fila (mock).
func (m *MockTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return nil
}
