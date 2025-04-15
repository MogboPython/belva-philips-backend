package database

import (
	"log"
	"os"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"
	"github.com/MogboPython/belvaphilips_backend/pkg/model"

	fiberLogger "github.com/gofiber/fiber/v2/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

// var DB_MIGRATOR gorm.Migrator

// ConnectDB connect to db
func ConnectDB() {
	var err error

	dsn := config.Config("DIRECT_URL")

	// Set logger level based on ENV
	logLevel := gormLogger.Silent
	if config.Config("ENV") == "development" {
		logLevel = gormLogger.Info
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.New(
			log.New(os.Stdout, "", 0),
			gormLogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		panic("failed to connect database")
	}

	// DB_MIGRATOR = DB.Migrator()

	fiberLogger.Info("Connection Opened to Database")
	DB.AutoMigrate(&model.User{}, &model.Order{})
	fiberLogger.Info("Database Migrated")

}
