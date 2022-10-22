package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
)

type Server struct{}

// Initialize initializes a web server
func (s *Server) Initialize() {
	mux := http.NewServeMux()
	mux.HandleFunc("/update", update)

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

func update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Successfully updated Planetscale Database"))

	env := "DEV"
	db, err := sql.Open("mysql", os.Getenv(env))
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
