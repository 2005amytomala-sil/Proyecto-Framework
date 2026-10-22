package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"multicatalogo-backend/interfaces"
	"multicatalogo-backend/models"
	"multicatalogo-backend/repository"
)

type ProdController struct {
	cache interfaces.CacheRepository
	repo  *repository.Repository
}

func NewProdController(cache interfaces.CacheRepository, repo *repository.Repository) *ProdController {
	return &ProdController{cache: cache, repo: repo}
}

// GetProductos implementa Cache-Aside con TTL de cinco minutos y tolera que Redis no esté disponible.
func (pc *ProdController) GetProductos(c *fiber.Ctx) error {
	ctx := context.Background()
	query := strings.TrimSpace(c.Query("q"))
	if query == "" && pc.cache != nil {
		if cached, err := pc.cache.Get(ctx, "productos:all"); err == nil {
			var productos []models.Producto
			if json.Unmarshal([]byte(cached), &productos) == nil {
				return c.JSON(productos)
			}
		}
	}
	if pc.repo == nil {
		return c.Status(503).JSON(models.APIError{Status: 503, Message: "Base de datos no disponible"})
	}
	var productos []models.Producto
	var err error
	if query == "" {
		productos, err = pc.repo.ObtenerProductos()
	} else {
		productos, err = pc.repo.BuscarProductos(query)
	}
	if err != nil {
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener productos"})
	}
	if query == "" && pc.cache != nil {
		if data, err := json.Marshal(productos); err == nil {
			_ = pc.cache.Set(ctx, "productos:all", data, 5*time.Minute)
		}
	}
	return c.JSON(productos)
}

func (pc *ProdController) GetProductoPorID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 1 {
		return c.Status(400).JSON(models.APIError{Status: 400, Message: "ID inválido: debe ser un número positivo", Details: map[string]string{"received": idParam}})
	}
	if pc.repo == nil {
		return c.Status(503).JSON(models.APIError{Status: 503, Message: "Base de datos no disponible"})
	}
	producto, err := pc.repo.ObtenerProductoPorID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return c.Status(404).JSON(models.APIError{Status: 404, Message: "Producto no encontrado"})
		}
		return c.Status(500).JSON(models.APIError{Status: 500, Message: "Error al obtener producto"})
	}
	return c.JSON(producto)
}
