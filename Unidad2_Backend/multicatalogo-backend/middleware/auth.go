// Declaramos el paquete middleware para los middlewares de autenticación.
package middleware

import (
	// Importamos strings para manipular strings.
	"strings"

	"github.com/gofiber/fiber/v2"

	// Importamos nuestras utilidades de JWT.
	"multicatalogo-backend/utils"
)

// RequireAuth es un middleware que requiere autenticación JWT.
func RequireAuth(c *fiber.Ctx) error {
	// Obtenemos el header Authorization.
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Authorization header requerido",
		})
	}

	// El header debe tener el formato "Bearer <token>".
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Formato de Authorization inválido. Debe ser 'Bearer <token>'",
		})
	}

	// Validamos el token.
	claims, err := utils.ValidateToken(parts[1])
	if err != nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "Token inválido o expirado",
			"details": err.Error(),
		})
	}

	// Guardamos los claims en el contexto de Fiber.
	c.Locals("email", claims.Email)
	c.Locals("rol", claims.Rol)

	// Continuamos con el siguiente handler.
	return c.Next()
}

// RequireAdmin es un middleware que requiere rol de admin.
func RequireAdmin(c *fiber.Ctx) error {
	// Obtenemos el rol del contexto (debe venir de RequireAuth).
	rol := c.Locals("rol")
	if rol == nil {
		return c.Status(401).JSON(fiber.Map{
			"status":  401,
			"message": "No autenticado",
		})
	}

	// Verificamos que el rol sea admin.
	if rol != "admin" {
		return c.Status(403).JSON(fiber.Map{
			"status":  403,
			"message": "Acceso denegado. Se requiere rol de admin",
		})
	}

	// Continuamos con el siguiente handler.
	return c.Next()
}
