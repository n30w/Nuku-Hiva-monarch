package main

import (
	"database/sql"
	"fmt"
	"os"
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
	fmt.Println("Successfully connected to PlanetScale!")

	posts, _ := ReadRedditData()

	for _, i := range posts.Rows {
		fmt.Printf("%d %s %s \n", i.Col1, i.Col2, i.Col3)
	}

	// for _, i := range comments.Rows {
	// 	fmt.Printf("%d %s %s \n", i.Col1, i.Col2, i.Col3)
	// }

	// err = psdb.InsertToSQL(posts)
	// if err != nil {
	// 	log.Fatal(err)
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
