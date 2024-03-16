package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	connStr := fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", getEnvVariable("DATABASE_USER"), getEnvVariable("DATABASE_NAME"), getEnvVariable("DATABASE_PASSWORD"))

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	return db, nil
}

func Setup(db *sql.DB) error {
	sqlFile, err := os.ReadFile("scheme.sql")
	if err != nil {
		log.Fatalf("Error reading SQL file: %v", err)
	}

	sqlFileContent := string(sqlFile)

	queries := strings.Split(sqlFileContent, ";")
	for _, query := range queries {
		var trimmedQuery = strings.TrimSpace(query)
		if trimmedQuery != "" {
			_, err := db.Exec(trimmedQuery)
			if err != nil {
				log.Fatalf("Error executing SQL query: %v", err)
			}
		}
	}
	return nil
}

func getEnvVariable(key string) string {
	envVariable, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Could not establish database connection: no %s env var", key)
	}

	return envVariable
}
