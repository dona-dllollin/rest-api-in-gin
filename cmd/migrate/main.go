package main

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please provide a migration direction: 'up' or 'down'")
	}

	direction := os.Args[1]

	// Database connection string
	dbURL := "postgres://postgres:postgres@localhost:5431/eventdb?sslmode=disable"

	// Create a new migrate instance
	m, err := migrate.New(
		"file://cmd/migrate/migrations", // Path to your migration files
		dbURL,
	)
	if err != nil {
		log.Fatalf("failed to create migrate instamce : %v", err)
	}

	switch direction {
	case "up":
		err = m.Up()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration up failder : %v", err)
		}
		log.Println("migration up completed")
	case "down":
		err = m.Down()
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("migration down failed: %v", err)
		}
		log.Println("migration down completed")

	default:
		log.Fatal("invalid direction. Use 'up' or 'down'")
	}

}
