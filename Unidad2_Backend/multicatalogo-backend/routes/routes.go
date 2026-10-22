// Declaramos el paquete routes, cuya única responsabilidad es registrar los endpoints de la API.
package routes

import (
	// Importamos Fiber para poder recibir y configurar la instancia principal de la aplicación.
	"github.com/gofiber/fiber/v2"
	// Importamos nuestros controladores.
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/middleware"
	// Importamos nuestros middlewares de autenticación.
)

// SetupRoutes es una función que recibe los controladores inyectados y configura las rutas.
func SetupRoutes(
	app *fiber.App,
	authController *controllers.AuthController,
	healthController *controllers.HealthController,
	prodController *controllers.ProdController,
	redController *controllers.RedController,
	checkoutController *controllers.CheckoutController,
	backofficeController *controllers.BackofficeController,
) {
	// Definimos una ruta HTTP GET en "/api/health" y la enlazamos al controlador.
	app.Get("/api/health", healthController.Health)

	// Definimos una ruta HTTP POST en "/api/login" y la enlazamos al controlador.
	app.Post("/api/login", authController.Login)

	// Definimos una ruta HTTP POST en "/api/register" y la enlazamos al controlador.
	app.Post("/api/register", authController.Register)

	// Definimos una ruta HTTP GET en "/api/productos" y la enlazamos al controlador.
	app.Get("/api/productos", prodController.GetProductos)

	// Definimos una ruta HTTP GET en "/api/productos/:id" y la enlazamos al controlador.
	app.Get("/api/productos/:id", prodController.GetProductoPorID)

	// Definimos una ruta HTTP GET en "/api/red" y la enlazamos al controlador.
	// Requiere autenticación JWT.
	app.Get("/api/red", middleware.RequireAuth, redController.GetRed)
	app.Get("/api/red/mia", middleware.RequireAuth, redController.GetRed)

	// Definimos una ruta HTTP POST en "/api/checkout" y la enlazamos al controlador.
	// Requiere autenticación JWT.
	app.Post("/api/checkout", middleware.RequireAuth, checkoutController.Checkout)

	admin := func(c *fiber.Ctx) error { return middleware.RequireAdmin(c) }
	app.Post("/api/productos", middleware.RequireAuth, admin, backofficeController.CreateProduct)
	app.Put("/api/productos/:id", middleware.RequireAuth, admin, backofficeController.UpdateProduct)
	app.Delete("/api/productos/:id", middleware.RequireAuth, admin, backofficeController.DeleteProduct)
	app.Get("/api/admin/productos", middleware.RequireAuth, admin, backofficeController.AdminProducts)
	app.Get("/api/admin/usuarios", middleware.RequireAuth, admin, backofficeController.AdminUsers)
	app.Get("/api/pedidos/mios", middleware.RequireAuth, backofficeController.MyOrders)
	app.Get("/api/pedidos", middleware.RequireAuth, admin, backofficeController.AdminOrders)
	app.Get("/api/pedidos/:id", middleware.RequireAuth, backofficeController.OrderDetail)
	app.Patch("/api/pedidos/:id/estado", middleware.RequireAuth, admin, backofficeController.UpdateOrderState)
	app.Get("/api/comisiones/mias", middleware.RequireAuth, backofficeController.MyCommissions)
	app.Post("/api/liquidaciones", middleware.RequireAuth, admin, backofficeController.Liquidate)
	app.Get("/api/liquidaciones", middleware.RequireAuth, admin, backofficeController.Liquidations)
	app.Get("/api/liquidaciones/mias", middleware.RequireAuth, backofficeController.MyLiquidations)
	app.Get("/api/reportes/resumen", middleware.RequireAuth, admin, backofficeController.Reports)
	app.Get("/api/reportes/top-productos", middleware.RequireAuth, admin, backofficeController.TopProducts)
	app.Get("/api/reportes/top-referidos", middleware.RequireAuth, admin, backofficeController.TopReferidos)
	app.Get("/api/reportes/comisiones-pendientes", middleware.RequireAuth, admin, backofficeController.PendingCommissions)
}
