// Declaramos el paquete controllers para agrupar las funciones que manejan la lógica de negocio de las rutas.
package controllers

import (
	// Importamos context para manejar el contexto de las operaciones.
	"context"
	// Importamos fmt para formatear strings.
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"

	// Importamos nuestros modelos.
	"multicatalogo-backend/models"
	// Importamos nuestro repositorio de checkout.
	"multicatalogo-backend/repository"
)

// CheckoutController es el controlador que maneja el proceso de checkout.
type CheckoutController struct {
	checkoutRepo *repository.CheckoutRepository
}

// NewCheckoutController crea una nueva instancia de CheckoutController.
func NewCheckoutController(checkoutRepo *repository.CheckoutRepository) *CheckoutController {
	return &CheckoutController{
		checkoutRepo: checkoutRepo,
	}
}

// Checkout procesa el pedido del cliente usando transacciones ACID y concurrencia.
func (cc *CheckoutController) Checkout(c *fiber.Ctx) error {
	ctx := context.Background()

	// Parsear el cuerpo de la petición.
	var req models.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return apiError(c, 400, "Cuerpo de petición inválido")
	}

	// Validar campos obligatorios.
	if req.MetodoPago == "" {
		return apiError(c, 400, "Método de pago es obligatorio")
	}
	if len(req.Items) == 0 {
		return apiError(c, 400, "El carrito no puede estar vacío")
	}

	// La identidad proviene del JWT; el email del body no es una fuente confiable.
	email, ok := c.Locals("email").(string)
	if !ok || email == "" {
		return c.Status(401).JSON(models.APIError{Status: 401, Message: "Sesión no válida"})
	}

	// CrearPedido usa los precios de PostgreSQL, bloquea stock y confirma todo junto.
	pedido, err := cc.checkoutRepo.CrearPedido(ctx, email, req.MetodoPago, req.Items)
	if err != nil {
		status := 500
		if strings.Contains(err.Error(), "no encontrado") || strings.Contains(err.Error(), "inválido") || strings.Contains(err.Error(), "vacío") {
			status = 400
		} else if strings.Contains(err.Error(), "stock insuficiente") {
			status = 409
		}
		return apiError(c, status, "No se pudo procesar el pedido", err.Error())
	}

	// Paso 2: Si el COMMIT fue exitoso, lanzamos una Goroutine para calcular comisiones.
	// Esto NO bloquea la respuesta al cliente.
	go func() {
		err := cc.checkoutRepo.CalcularComisiones(context.Background(), email, pedido.Total)
		if err != nil {
			fmt.Printf("Error al calcular comisiones (Goroutine): %v\n", err)
		} else {
			fmt.Printf("Comisiones calculadas exitosamente para %s: %.2f\n", email, pedido.Total)
		}
	}()

	// Paso 3: Retornamos respuesta inmediata al cliente (no esperamos a las comisiones).
	return c.Status(201).JSON(fiber.Map{
		"message": "Pedido procesado exitosamente",
		"pedido":  pedido,
	})
}
