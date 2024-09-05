package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"os"
	"strconv"
)

func NewSQLX() *sqlx.DB {
	username := os.Getenv("DATABASE_USERNAME")
	password := os.Getenv("DATABASE_PASSWORD")
	host := os.Getenv("DATABASE_HOST")
	port := os.Getenv("DATABASE_PORT")
	database := os.Getenv("DATABASE_NAME")

	portInt, err := strconv.Atoi(port)
	if err != nil {
		panic("failed to convert DATABASE_PORT to int: " + err.Error())
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, portInt, database)

	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return db
}
