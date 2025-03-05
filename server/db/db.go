package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
)

var DB *sql.DB

func InitDB() {
	connectionString := os.Getenv("DATABASE_URL")
	if connectionString == "" {
		connectionString = "postgres://historymapuser:your_strong_password@localhost:5432/historical_maps?sslmode=disable"
	}
	var err error
	DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal("Error connecting to the database: ", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatal("Error pinging the database: ", err)
	}
	fmt.Println("Successfully connected to the database")

	createTables()
}

func GetDB() *sql.DB {
	return DB
}

func createTables() {
	createMapTableSQL := `
  CREATE TABLE IF NOT EXISTS historical_maps(
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year INTEGER,
    image_path  VARCHAR(255),
    bounds GEOMETRY(Polygon, 4326)
  );
  `

	_, err := DB.Exec(createMapTableSQL)
	if err != nil {
		log.Fatal("Error creating historical_maps table:", err)
	}
	fmt.Println("historical_maps table created (or already exists)")
}
