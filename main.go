package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	db, err := sql.Open("mysql", os.Getenv("DEV"))
	if err != nil {
		panic(err)
	}

	defer db.Close()

	if err := db.Ping(); err != nil {
		panic(err)
	}

	key := &Key{}
	redditPosts := &Table[Row[id, text]]{Name: "posts"}
	dbPosts := &Table[Row[id, text]]{Name: "posts"}
	redditComments := &Table[Row[id, text]]{Name: "comments"}
	dbComments := &Table[Row[id, text]]{Name: "comments"}

	GrabSaved(redditPosts, redditComments, key, 1)

	psdb := &PlanetscaleDB{db}
	err = psdb.RetrieveSQL(dbPosts)
	if err != nil {
		log.Fatalln(err)
	}
	err = psdb.RetrieveSQL(dbComments)
	if err != nil {
		log.Fatalln(err)
	}

	psdb.UpdateSQL(dbPosts, redditPosts)
	psdb.UpdateSQL(dbComments, redditComments)
}
