package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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

func InsertToSQL(db *sql.DB, table *Table[Row[id, text]]) error {

	var query string
	var inserts []string
	var params []interface{}
	insertion := "(?, ?, ?, ?, ?)"

	switch table.Name {
	case "posts":
		query = "INSERT INTO posts (id, name, url, subreddit, media_url) VALUES "
	case "comments":
		query = "INSERT INTO comments (id, author, body, url, subreddit) VALUES "
	default:
		errors.New("Not a valid table name!")
	}

	for _, v := range table.Rows {
		inserts = append(inserts, insertion)
		params = append(
			params,
			v.Col1,
			v.Col2,
			v.Col3,
			v.Col4,
			v.Col5,
		)
	}

	qv := strings.Join(inserts, ",")

	query += qv

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	statement, err := db.PrepareContext(ctx, query)

	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	defer statement.Close()

	res, err := statement.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	log.Printf("%d post rows created ", rows)
	fmt.Printf("%d post rows created", rows)

	return nil
}
