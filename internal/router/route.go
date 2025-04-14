package router

import (
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/middleware"
	swagger "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, orderHandler *handler.OrderHandler) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	swaggerCfg := swagger.Config{
		BasePath: "/cmd/app/docs",
		FilePath: "./cmd/app/docs/swagger.json",
	}

	app.Use(swagger.New(swaggerCfg))

	// grouping
	api := app.Group("/api/v1")
	api.Post("/admin/login", adminHandler.AdminLogin)
	api.Post("/contact", handler.ContactUs)
	{
		user := api.Group("/users")
		user.Post("/token", userHandler.CreateUserAccessToken)
		user.Get("/get_user", middleware.Protected(), userHandler.GetUserByEmail)
		user.Post("/", middleware.Protected(), userHandler.CreateUser)
		user.Put("/:id/membership", userHandler.UpdateMembershipStatus)
		// v1.Get("/:id", middleware.Protected(), userHandler.GetUserByID)
		// v1.Put("/:id", handler.UpdateUser)
	}
	{
		admin := api.Group("/admin", middleware.Protected(), middleware.AdminRole())
		admin.Get("/get_users", adminHandler.GetAllUsers)
		admin.Get("/user/:id", adminHandler.GetUserByID)
	}
	{
		order := api.Group("/orders/", middleware.Protected())
		order.Post("/", orderHandler.CreateOrder)
		order.Get("/", middleware.AdminRole(), orderHandler.GetAllOrders)
		order.Get("/user/:userId", orderHandler.GetOrdersByUserID)
		order.Get("/:id", orderHandler.GetOrderByID)
		order.Put("/:id/status", middleware.AdminRole(), orderHandler.UpdateOrderStatus)
	}

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(404) // => 404 "Not Found"
	})
}
