package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

type RedController struct {
	cache interfaces.CacheRepository
	repo  *repository.Repository
}

func NewRedController(cache interfaces.CacheRepository, repo *repository.Repository) *RedController {
	return &RedController{cache: cache, repo: repo}
}

// GetRed only exposes the authenticated user's own referral subtree.
func (rc *RedController) GetRed(c *fiber.Ctx) error {
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Sesión no válida"})
	}
	ctx := context.Background()
	cacheKey := fmt.Sprintf("red:tree:%s", email)
	if rc.cache != nil {
		if cached, err := rc.cache.Get(ctx, cacheKey); err == nil {
			var tree models.NodoRed
			if json.Unmarshal([]byte(cached), &tree) == nil {
				return c.JSON(tree)
			}
		}
	}
	if rc.repo == nil {
		return c.Status(503).JSON(models.APIError{Status: 503, Message: "Base de datos no disponible"})
	}
	tree, err := rc.repo.ObtenerRed(email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Red no encontrada"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener red"})
	}
	if rc.cache != nil {
		if data, err := json.Marshal(tree); err == nil {
			_ = rc.cache.Set(ctx, cacheKey, data, 10*time.Minute)
		}
	}
	return c.JSON(tree)
}
