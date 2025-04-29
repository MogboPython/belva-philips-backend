package main

import (
	_ "github.com/MogboPython/belvaphilips_backend/cmd/app/docs"
	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/internal/database"
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/internal/router"
	"github.com/MogboPython/belvaphilips_backend/internal/service"
	_ "github.com/lib/pq"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

// @title						Belva Philips Backend API
// @version					    1.0
// @description				    This is an backend API for Belva Philips website
//
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "UTC",
	}))
	app.Use(cors.New())

	if err := database.ConnectDB(); err != nil {
		log.Errorf("Failed to connect to the database: %v", err)
	}

	if err := database.MigrateDB(); err != nil {
		log.Errorf("Failed to migrate database: %v", err)
	}

	db := database.DB
	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	postRepo := repository.NewPostRepository(db)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	adminService := service.NewAdminService(userRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	orderService := service.NewOrderService(orderRepo)
	orderHandler := handler.NewOrderHandler(orderService)

	postService := service.NewPostService(postRepo)
	postHandler := handler.NewPostHandler(postService)

	app.Get("/swagger/*", swagger.HandlerDefault)

	router.SetupRoutes(app, userHandler, adminHandler, orderHandler, postHandler)

	if err := app.Listen(":" + config.Config("PORT")); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
