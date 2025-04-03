package main

import (
	"fmt"

	_ "github.com/MogboPython/belvaphilips_backend/cmd/app/docs"
	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/internal/database"
	"github.com/MogboPython/belvaphilips_backend/internal/handler"
	"github.com/MogboPython/belvaphilips_backend/internal/repository"
	"github.com/MogboPython/belvaphilips_backend/internal/router"
	"github.com/MogboPython/belvaphilips_backend/internal/service"

	// "github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// @title	Belva Philips Backend API
// @version		1.0
// @description	This is an backend Api just for Belva Philips website
// @termsOfService	http://swagger.io/terms/
// @host localhost:8080
// @BasePath /
func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	database.ConnectDB()
	db := database.DB

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)

	router.SetupRoutes(app, userHandler)

	app.Listen(fmt.Sprintf(":%s", config.Config("PORT")))
}
