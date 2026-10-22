// Declaramos el paquete config para agrupar las configuraciones de la aplicación.
package config

import (
	"bufio"
	"context"
	// Importamos el paquete database/sql para manejar conexiones a bases de datos.
	"database/sql"
	// Importamos el paquete fmt para formatear strings.
	"fmt"
	"os"
	"strings"
	"time"

	// Importamos el driver de PostgreSQL.
	_ "github.com/lib/pq"
)

// DB es la variable global que contiene la conexión a PostgreSQL.
var DB *sql.DB

// LoadEnv carga un .env local sin reemplazar variables ya definidas por el entorno.
func LoadEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key, value = strings.TrimSpace(key), strings.Trim(strings.TrimSpace(value), "\"'")
		if key != "" && os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}
}

// InitPostgres inicializa la conexión con PostgreSQL usando variables de entorno.
func InitPostgres() {
	// Leemos las configuraciones desde variables de entorno, con fallbacks.
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5434"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "multicatalogo"
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		dbPassword = "multicatalogo123"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "multicatalogo"
	}

	// Construimos el string de conexión (DSN).
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Abrimos la conexión a la base de datos.
	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		panic("Error al abrir la conexión a PostgreSQL: " + err.Error())
	}

	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(5 * time.Minute)
	// Verificamos la conexión con un límite de tiempo.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = DB.PingContext(ctx)
	if err != nil {
		fmt.Println("No se pudo conectar con PostgreSQL:", err.Error())
		_ = DB.Close()
		DB = nil
		return
	}
}
