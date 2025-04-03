package router

import (
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/middleware"
	swagger "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler) {
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
	v1 := api.Group("/user")

	// routes
	v1.Post("/token", userHandler.CreateUserAccessToken)
	v1.Get("/", middleware.Protected(), userHandler.GetAllUsers)
	v1.Get("/:id", middleware.Protected(), userHandler.GetUserByID)
	v1.Post("/", middleware.Protected(), userHandler.CreateUser)
	// v1.Put("/:id", handler.UpdateUser)

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}
