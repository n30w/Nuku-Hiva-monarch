package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // The underscore on imports autoloads the dependency. Do not need to call something like "godotenv.Load()"
	_ "github.com/joho/godotenv/autoload"
)

func GetDatabase() (*sql.DB, error) {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	return db, err
}

func Connect() {
	db, err := GetDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PlanetScale!")

	var title string
	id := 2
	err = db.QueryRow("SELECT name FROM posts WHERE id = ?", id).Scan(&title)
	if err == sql.ErrNoRows {
		log.Fatal("no rows returned")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
}

func retrieveBy[T Col](db *sql.DB, table *Table, retrieval, want *T) T {
	q := "SELECT " + string(retrieval) + " FROM " + table + " WHERE " + want

	err := db.QueryRow(q).Scan(&retrieval)
	if err == sql.ErrNoRows {
		log.Fatal("no rows returned")
	} else if err != nil {
		log.Fatal(err)
	}
	return *retrieval
}
