// Declaramos el paquete config para agrupar las configuraciones de la aplicación.
package config

import (
	// Importamos el paquete context para manejar timeouts y cancelaciones.
	"context"
	"fmt"
	// Importamos el paquete os para leer variables de entorno.
	"os"
	"time"
	// Importamos el cliente oficial de Redis para Go.
	"github.com/redis/go-redis/v9"
)

// RedisClient es la variable global que contiene la conexión a Redis.
var RedisClient *redis.Client

// InitRedis inicializa la conexión con Redis usando variables de entorno.
func InitRedis() {
	// Leemos el host de Redis desde la variable de entorno, con fallback a localhost.
	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	// Leemos el puerto de Redis desde la variable de entorno, con fallback a 6379.
	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	// Inicializamos el cliente de Redis con la configuración.
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisHost + ":" + redisPort,
		Password: "", // Sin contraseña por defecto
		DB:       0,  // Base de datos 0
	})

	// Verificamos la conexión con un ping.
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		fmt.Println("Redis no disponible; se atenderán las lecturas desde PostgreSQL:", err.Error())
		_ = RedisClient.Close()
		RedisClient = nil
	}
}
