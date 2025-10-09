package config

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading environment variables directly")
	}
}

func ConnectDB() *sql.DB {
	dsn := os.Getenv("SUPABASE_URL")
	if dsn == "" {
		log.Fatal("SUPABASE_URL not set in environment")
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("Successfully connected to Supabase!")
	return db
}
