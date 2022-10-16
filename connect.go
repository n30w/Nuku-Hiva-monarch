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

func Connect(db *sql.DB) {
	var title string
	id := 2
	err := db.QueryRow("SELECT name FROM posts WHERE id = ?", id).Scan(&title)
	if err == sql.ErrNoRows {
		log.Fatal("no rows returned")
	} else if err != nil {
		log.Fatal(err)
	}
	fmt.Println(title)
}

// https://golangbot.com/mysql-create-table-insert-row/
func UploadDataToPlanetScale(db *sql.DB, p []*Post, c []*Comment) {
	s := ""
	result, err := db.Exec(s)
}

// // Insert data into table
// func Insert[T Col](db *sql.DB, table *Table) error {
// 	q := "INSERT INTO " + table.Name + " "

// 	switch table.Name {
// 	case "posts":
// 		q += "posts"
// 		ferry := &Post{}
// 	case "comments":
// 		q += "comments"
// 		ferry := &Comment{}
// 	default:
// 		return errors.New("Not a valid table!")
// 	}

// 	result, err := db.Exec("INSERT INTO customers (name) VALUES (?)", "Alice")

// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return nil
// }
