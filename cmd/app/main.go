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

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)

//	@title						Belva Philips Backend API
//	@version					1.0
//	@description				This is an backend API for Belva Philips website
//	@termsOfService				http://swagger.io/terms/
//	@host						localhost:8080
//	@BasePath					/
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func main() {
	app := fiber.New()
	app.Use(logger.New(logger.Config{
		Format:     "${cyan}[${time}] ${white}${pid} ${red}${status} ${blue}[${method}] ${white}${path}\n",
		TimeFormat: "02-Jan-2006",
		TimeZone:   "UTC",
	}))
	app.Use(cors.New())

	database.ConnectDB()
	db := database.DB

	// Initialize repository, service, and handler
	userRepo := repository.NewUserRepository(db)

	userService := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userService)
	adminService := service.NewAdminService(userRepo)
	adminHandler := handler.NewAdminHandler(adminService)

	app.Get("/swagger/*", swagger.HandlerDefault)

	router.SetupRoutes(app, userHandler, adminHandler)

	app.Listen(fmt.Sprintf(":%s", config.Config("PORT")))
}
