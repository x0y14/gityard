package database

import (
	"fmt"
	"gityard-api/config"
	"log"
	"strconv"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDB() {
	var err error
	p := config.Config("DB_PORT")
	port, err := strconv.ParseUint(p, 10, 32)
	if err != nil {
		panic("failed to parse database port")
	}

	utc, err := time.LoadLocation("UTC")
	if err != nil {
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
	log.Println(dsn)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	fmt.Println("Connection Opened to Database")
	//err = DB.AutoMigrate(
	//	&model.User{},
	//	&model.UserCredential{},
	//	&model.UserRefreshToken{},
	//	&model.Handlename{},
	//	&model.Account{},
	//	&model.AccountProfile{},
	//	&model.AccountPublicKey{},
	//	&model.Repository{},
	//)
	//if err != nil {
	//	slog.Error("failed to migrate db", "detail", err)
	//	panic("failed to migrate db")
	//}
	//fmt.Println("Database Migrated")
}
