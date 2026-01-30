package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"database/sql"

	_ "github.com/lib/pq"
)

type RootConfig struct {
	DBUrl string
}

func NewPostgres() (*sql.DB, error) {
	cfg := LoadConfig()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxLifetime(time.Hour)

	log.Println("Postgres connected and configured")
	return db, nil
}

func LoadConfig() *RootConfig {
	err := godotenv.Load() // Fails to find file in cloud
	if err != nil {
		// This log line executes in the cloud
		log.Println("Error loading .env file: Falling back to environment variables.")
		log.Println(err)
	}

	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("DATABASE_URL is required")
	}

	return &RootConfig{
		DBUrl: dbUrl,
	}
}
