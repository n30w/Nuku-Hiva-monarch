package main

import (
	"database/sql"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// psdb := &PlanetscaleDB{db}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	posts, comments := ReadAllRedditSaved()

	posts.List()
	comments.List()

	// {
	// 	err = psdb.InsertToSQL(posts)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}

	// 	err = psdb.InsertToSQL(comments)
	// 	if err != nil {
	// 		log.Fatalln(err)
	// 	}
	// }
	// API FUNCTIONS
	//
	// MASS REFRESH:
	// 1) Get data from Reddit
	// 2) Store data into structs
	// 3) Transfer struct data over to PlanetScale SQL database
	//
	// CONSISTENT UPDATES:
	// 1) Get MOST RECENT saved posts from Reddit every 24 hours
	// 2) Store data into structs
	// 3) Transfer struct data over to PlanetScale SQL database
}
