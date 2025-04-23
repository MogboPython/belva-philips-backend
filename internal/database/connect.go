package database

import (
	"log"
	"os"
	"time"

	"github.com/MogboPython/belvaphilips_backend/internal/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// ConnectDB connect to db
func ConnectDB() error {
	var err error

	dsn := config.Config("DIRECT_URL")

	// Set logger level based on ENV
	logLevel := gormlogger.Silent
	if config.Config("ENV") == "development" {
		logLevel = gormlogger.Info
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.New(
			log.New(os.Stdout, "", 0),
			gormlogger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logLevel,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		return err
	}

	SQLDb, _ = DB.DB()

	return nil
}

// MigrateDB to migrate db
func MigrateDB() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(SQLDb, "internal/database/migrations"); err != nil && err != goose.ErrAlreadyApplied {
		return err
	}

	log.Println("Database Migrated")

	return nil
}
