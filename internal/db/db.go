package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", GetEnvVariable("DATABASE_USER"), GetEnvVariable("DATABASE_NAME"), GetEnvVariable("DATABASE_PASSWORD"))

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return DB, nil
}

func GetEnvVariable(key string) string {
	envVariable, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Could not establish database connection: no %s env var", key)
	}

	return envVariable
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("Database connection closed")
	}
}
