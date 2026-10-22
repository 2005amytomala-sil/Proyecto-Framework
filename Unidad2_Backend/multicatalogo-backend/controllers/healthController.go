package controllers

import (
	"context"
	"database/sql"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type HealthController struct {
	db    *sql.DB
	redis *redis.Client
}

func NewHealthController(db *sql.DB, redisClient *redis.Client) *HealthController {
	return &HealthController{db: db, redis: redisClient}
}

func (hc *HealthController) Health(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if hc.db == nil || hc.db.PingContext(ctx) != nil {
		return c.Status(503).JSON(fiber.Map{"status": "error", "message": "PostgreSQL no disponible"})
	}
	redisStatus := "ok"
	if hc.redis == nil || hc.redis.Ping(ctx).Err() != nil {
		redisStatus = "degraded"
	}
	return c.JSON(fiber.Map{"status": "ok", "message": "Servidor Go/Fiber operativo", "postgres": "ok", "redis": redisStatus})
}
