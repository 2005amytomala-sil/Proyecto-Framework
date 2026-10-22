// Declaramos el paquete repository para las implementaciones concretas de los repositorios.
package repository

import (
	// Importamos context para manejar el contexto de las operaciones.
	"context"
	"errors"
	// Importamos time para definir la duración de expiración.
	"time"

	// Importamos el cliente de Redis.
	"github.com/redis/go-redis/v9"

	// Importamos nuestras interfaces.
	"multicatalogo-backend/interfaces"
)

// RedisCache es la implementación concreta de CacheRepository usando Redis.
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache crea una nueva instancia de RedisCache.
func NewRedisCache(client *redis.Client) interfaces.CacheRepository {
	return &RedisCache{client: client}
}

// Get obtiene un valor del caché por su clave.
func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	if r.client == nil {
		return "", errors.New("redis no disponible")
	}
	return r.client.Get(ctx, key).Result()
}

// Set guarda un valor en el caché con un tiempo de expiración.
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration interface{}) error {
	if r.client == nil {
		return errors.New("redis no disponible")
	}
	var exp time.Duration
	switch v := expiration.(type) {
	case time.Duration:
		exp = v
	case int:
		exp = time.Duration(v) * time.Second
	default:
		exp = 0
	}
	return r.client.Set(ctx, key, value, exp).Err()
}

// Delete elimina una clave del caché.
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	if r.client == nil {
		return errors.New("redis no disponible")
	}
	return r.client.Del(ctx, key).Err()
}
