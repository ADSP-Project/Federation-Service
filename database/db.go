package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func DbConn() (db *sql.DB) {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	dbDriver := "postgres"
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dbInfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", dbUser, dbPass, dbName)

	db, err = sql.Open(dbDriver, dbInfo)
	if err != nil {
		panic(err.Error())
	}
	
	return db
}

func GetWebhookURL(shopId string) (string, error) {
	db := DbConn()

	var webhookURL string
	err := db.QueryRow("SELECT webhookurl FROM shops WHERE id = $1", shopId).Scan(&webhookURL)

	db.Close()

	if err != nil {
		if err == sql.ErrNoRows {
			// No rows were returned - handle according to your requirements
			log.Printf("No shop found with id: %s\n", shopId)
		}
		return "", err
	}

	return webhookURL, nil
}