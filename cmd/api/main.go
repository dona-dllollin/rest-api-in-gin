package main

import (
	"database/sql"
	"log"

	_ "rest-api-in-gin/docs"
	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

// @title Go Gin Rest API
// @version 1.0
// @description A rest APPI in Go using Gin framework
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format **Bearer &lt:token&gt;**

type application struct {
	port      int
	jwtSecret string
	models    database.Models
	redis     *redis.Client
}

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5431/eventdb?sslmode=disable")

	if err != nil {
		log.Fatal(err)
	}

	// Test koneksi
	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	log.Println("Database connection Succesfuly")

	defer db.Close()

	models := database.NewModels(db)
	app := &application{
		port:      env.GetEnvInt("PORT", 8080),
		jwtSecret: env.GetEnvString("JWT_SECRET", "some-secret-123456"),
		models:    models,
		redis:     redis.NewClient(&redis.Options{Addr: "localhost:6379"}),
	}

	if err := app.serve(); err != nil {
		log.Fatal(err)
	}

}
