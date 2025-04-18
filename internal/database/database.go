package database

import (
	"database/sql"

	"gorm.io/gorm"
)

// DB gorm connector
var DB *gorm.DB

var SQLDb *sql.DB
