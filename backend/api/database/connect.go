package database

import (
	"fmt"
	"gityard-api/config"
	"log/slog"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() {
	slog.Info("try connect to database")

	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		slog.Error("failed to parse database port", "detail", err)
		panic("failed to parse database port")
	}

	utc, err := time.LoadLocation("UTC")
	if err != nil {
		slog.Error("failed to load utc tz", "detail", err)
		panic("failed to load utc tz")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=%s",
		config.Config("DB_USER"),
		config.Config("DB_PASSWORD"),
		config.Config("DB_HOST"),
		port,
		config.Config("DB_NAME"),
		utc,
	)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("failed to connect database", "detail", err)
		panic("failed to connect database")
	}

	slog.Info("connection opened to database")
}
