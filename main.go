package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

const (
	PleasePopulateIDs = false
	version           = "1.0.3"
)

var (
	env = os.Getenv("ENVIRONMENT")
	db  *sql.DB
)

func main() {
	var err error
	db, err = sql.Open("mysql", os.Getenv(env))
	if err != nil {
		panic(Warn.Sprint(err))
	}

	if err := db.Ping(); err != nil {
		panic(Warn.Sprint(err))
	}

	server := &Server{
		RedditPosts:    &Table[Row[id, text]]{Name: "posts"},
		RedditComments: &Table[Row[id, text]]{Name: "comments"},
		DBPosts:        &Table[Row[id, text]]{Name: "posts"},
		DBComments:     &Table[Row[id, text]]{Name: "comments"},
		Key:            &Key{},
		PlanetscaleDB:  &PlanetscaleDB{db},
	}

	log.Print(Start.Sprintf("Starting andthensome %s %s", version, env))
	log.Print(Start.Sprint("Server listening on :4000"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.UpdateHandler)
	mux.HandleFunc("/areyouawake", server.AwakeHandler)

	// Only allow certain requests in Development environment only
	if env == "DEV" {
		mux.HandleFunc("/populate", server.PopulateHandler)
		mux.HandleFunc("/delete", server.ClearTableHandler(server)) // Why?
	}

	if err := http.ListenAndServe(":4000", mux); err != nil {
		fmt.Print(err)
		log.Fatal(Warn.Sprint(err))
	}
}
