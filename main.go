package main

import (
	"database/sql"
	"fmt"
	"os"
)

var db *sql.DB

func init() {
	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		panic(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to PlanetScale!")
}

func main() {

	x, _ := ReadRedditData()

	for _, a := range x {
		fmt.Printf("%d    %s\n", a.Id, a.Name)
	}
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
