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

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})

	swaggerCfg := swagger.Config{
		BasePath: "../../cmd/app/docs", // swagger ui base path
		FilePath: "../../cmd/app/docs/swagger.json",
		Path:     "swagger",
		Title:    "Belva Philips Backend API",
	}

	app.Use(swagger.New(swaggerCfg))
	// api := app.Group("/api")
	// v1 := api.Group("/v1", func(c *fiber.Ctx) error {
	// 	c.JSON(fiber.Map{
	// 		"message": "🐣 v1",
	// 	})
	// 	return c.Next()
	// })

	// grouping
	api := app.Group("/api/v1")
	v1 := api.Group("/user")

	// routes
	v1.Post("/token", userHandler.CreateUserAccessToken)
	// v1.Get("/", handler.GetAllUsers)
	// v1.Get("/:id", handler.GetSingleUser)
	v1.Post("/", middleware.Protected(), handler.CreateUser)
	// v1.Put("/:id", handler.UpdateUser)

}
