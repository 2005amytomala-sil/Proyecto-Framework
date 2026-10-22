package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	"multicatalogo-backend/config"
	"multicatalogo-backend/controllers"
	"multicatalogo-backend/repository"
	"multicatalogo-backend/routes"
	"multicatalogo-backend/utils"
)

func main() {
	config.LoadEnv()
	if err := utils.ConfigureJWTSecret(); err != nil {
		log.Fatal(err)
	}
	config.InitPostgres()
	if config.DB == nil {
		log.Fatal("PostgreSQL no está disponible; no se inicia el backend con datos simulados")
	}
	defer config.DB.Close()

	config.InitRedis() // Redis es opcional; la API sigue consultando PostgreSQL si falla.
	if config.RedisClient != nil {
		defer config.RedisClient.Close()
	}

	cacheRepo := repository.NewRedisCache(config.RedisClient)
	dbRepo := repository.NewPostgresDB(config.DB)
	dataRepo := repository.NewRepository(config.DB)
	checkoutRepo := repository.NewCheckoutRepository(dbRepo, cacheRepo)

	authController := controllers.NewAuthController(dataRepo)
	prodController := controllers.NewProdController(cacheRepo, dataRepo)
	redController := controllers.NewRedController(cacheRepo, dataRepo)
	healthController := controllers.NewHealthController(config.DB, config.RedisClient)
	checkoutController := controllers.NewCheckoutController(checkoutRepo)
	backofficeController := controllers.NewBackofficeController(config.DB, cacheRepo)

	app := fiber.New()
	allowedOrigins := os.Getenv("FRONTEND_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:5173,http://localhost:5174,http://localhost:5175,http://localhost:4173,http://127.0.0.1:5173,http://127.0.0.1:4173"
	}
	app.Use(cors.New(cors.Config{
		AllowOrigins: allowedOrigins,
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))
	routes.SetupRoutes(app, authController, healthController, prodController, redController, checkoutController, backofficeController)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
