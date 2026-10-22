// Declaramos el paquete repository para las implementaciones concretas de los repositorios.
package repository

import (
	// Importamos database/sql para manejar conexiones a bases de datos.
	"database/sql"
	"fmt"

	// Importamos nuestras interfaces.
	"multicatalogo-backend/interfaces"
)

// PostgresDB es la implementación concreta de DBRepository usando PostgreSQL.
type PostgresDB struct {
	db *sql.DB
}

// NewPostgresDB crea una nueva instancia de PostgresDB.
func NewPostgresDB(db *sql.DB) interfaces.DBRepository {
	return &PostgresDB{db: db}
}

// Begin inicia una transacción.
func (p *PostgresDB) Begin() (interfaces.Transaction, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return nil, err
	}
	return &PostgresTransaction{tx: tx}, nil
}

// Exec ejecuta una consulta SQL sin retornar filas.
func (p *PostgresDB) Exec(query string, args ...interface{}) error {
	_, err := p.db.Exec(query, args...)
	return err
}

// Query ejecuta una consulta SQL y retorna filas.
func (p *PostgresDB) Query(query string, args ...interface{}) ([]map[string]interface{}, error) {
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			row[col] = values[i]
		}
		results = append(results, row)
	}

	return results, nil
}

// QueryRow ejecuta una consulta SQL y retorna una sola fila.
func (p *PostgresDB) QueryRow(query string, args ...interface{}) (map[string]interface{}, error) {
	// Para simplificar, usamos Query y retornamos la primera fila.
	rows, err := p.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	if !rows.Next() {
		return nil, fmt.Errorf("no rows found")
	}

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	if err := rows.Scan(valuePtrs...); err != nil {
		return nil, err
	}

	result := make(map[string]interface{})
	for i, col := range columns {
		result[col] = values[i]
	}

	return result, nil
}

// PostgresTransaction es la implementación concreta de Transaction usando PostgreSQL.
type PostgresTransaction struct {
	tx *sql.Tx
}

// Commit confirma la transacción.
func (p *PostgresTransaction) Commit() error {
	return p.tx.Commit()
}

// Rollback revierte la transacción.
func (p *PostgresTransaction) Rollback() error {
	return p.tx.Rollback()
}

// Exec ejecuta una consulta SQL dentro de la transacción.
func (p *PostgresTransaction) Exec(query string, args ...interface{}) error {
	_, err := p.tx.Exec(query, args...)
	return err
}

// QueryRow ejecuta una consulta SQL dentro de la transacción y retorna una sola fila.
func (p *PostgresTransaction) QueryRow(query string, args ...interface{}) *sql.Row {
	return p.tx.QueryRow(query, args...)
}
