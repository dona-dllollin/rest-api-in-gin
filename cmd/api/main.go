package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "rest-api-in-gin/docs"
	"rest-api-in-gin/internal/database"
	"rest-api-in-gin/internal/env"

	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	//prometheus handler
	http.Handle("/metrics", promhttp.Handler())
	promServer := &http.Server{
		Addr: ":8081",
	}

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
	srv := app.serve()

	ctx, cancel := context.WithCancel(context.Background())
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	signal.Notify(sig, syscall.SIGTERM)

	go func() {
		<-sig
		slog.Info("Shutting down application")
		cancel()
	}()

	// start prometheus exporter
	go func() {
		slog.Info("Listening and serving prometheus exporter", slog.Int("port", 8081))
		if err := promServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed listening prometheus exporter", slog.Any("err", err))
			panic(err)
		}
	}()

	// start echo server
	go func() {
		slog.Info("Listening and serving HTTP", slog.Int("port", 8080))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed starting HTTP server", slog.Any("err", err))
			panic(err)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", slog.Any("err", err))
		panic(err)
	}
	if err := promServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Prometheus shutdown failed", slog.Any("err", err))
		panic(err)
	}

}
