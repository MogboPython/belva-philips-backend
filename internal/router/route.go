package router

import (
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/middleware"
	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

// SetupRoutes func
func SetupRoutes(app *fiber.App, userHandler *handler.UserHandler, adminHandler *handler.AdminHandler, orderHandler *handler.OrderHandler, postHandler *handler.PostHandler) {
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
	api.Post("/token", userHandler.CreateUserAccessToken)
	{
		user := api.Group("/users", middleware.Protected())
		user.Get("/:id", userHandler.GetUserByID)
		user.Post("/", userHandler.CreateUser)
		user.Put("/:id/membership", userHandler.UpdateMembershipStatus)
	}
	{
		admin := api.Group("/admin", middleware.Protected(), middleware.AdminRole())
		admin.Get("/get_users", adminHandler.GetAllUsers)
	}
	{
		order := api.Group("/orders/", middleware.Protected())
		// User-specific routes
		order.Get("/user/:userId", orderHandler.GetOrdersByUserID)

		// Admin-specific routes
		order.Get("/", middleware.AdminRole(), orderHandler.GetAllOrders)
		order.Put("/:order_id/status", middleware.AdminRole(), orderHandler.UpdateOrderStatus)

		// General routes
		order.Post("/", orderHandler.CreateOrder)
		order.Get("/:id", orderHandler.GetOrderByID)
	}
	{
		post := api.Group("/posts/")
		post.Post("/upload-image", middleware.Protected(), middleware.AdminRole(), postHandler.UploadImage)
		post.Post("/", middleware.Protected(), middleware.AdminRole(), postHandler.CreatePost)
		post.Get("/drafts", middleware.Protected(), middleware.AdminRole(), postHandler.GetAllDraftPosts)

		post.Get("/", postHandler.GetAllPosts)
		post.Get("/:id", postHandler.GetPostByID)
		post.Delete("/:id", middleware.Protected(), middleware.AdminRole(), postHandler.DeletePost)
	}

	// handle unavailable route
	app.Use(func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusNotFound)
	})
}

// order.Put("/:id/status", middleware.AdminRole(), orderHandler.UpdateOrderStatus)
// v1.Get("/:id", middleware.Protected(), userHandler.GetUserByID)
// v1.Put("/:id", handler.UpdateUser)
