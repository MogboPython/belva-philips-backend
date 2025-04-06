package router

import (
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/middleware"
	swagger "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, adminHandler *handler.AdminHandler) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	swaggerCfg := swagger.Config{
		BasePath: "/cmd/app/docs", // swagger ui base path
		FilePath: "./cmd/app/docs/swagger.json",
		// Path:     "swagger",
		// Title:    "Belva Philips Backend API",
	}

	app.Use(swagger.New(swaggerCfg))

	// grouping
	api := app.Group("/api/v1")
	{
		v1 := api.Group("/user")
		v1.Post("/token", userHandler.CreateUserAccessToken)
		v1.Get("/get_user", middleware.Protected(), userHandler.GetUserByEmail)
		v1.Post("/", middleware.Protected(), userHandler.CreateUser)
		// v1.Get("/:id", middleware.Protected(), userHandler.GetUserByID)
		// v1.Put("/:id", handler.UpdateUser)
	}
	{
		admin := api.Group("/admin")
		admin.Get("/login", middleware.AdminProtected(), adminHandler.AdminLogin)
		admin.Get("/:id", middleware.AdminProtected(), adminHandler.GetUserByID)
		admin.Get("/get_users", middleware.AdminProtected(), adminHandler.GetAllUsers)
	}

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}
