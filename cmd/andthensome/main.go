package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/n30w/andthensome/internal/style"

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
		panic(style.Warn.Sprint(err))
	}

	if err := db.Ping(); err != nil {
		panic(style.Warn.Sprint(err))
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
	mux.HandleFunc("/api/update", server.UpdateHandler)
	mux.HandleFunc("/api/scananddelete", server.ScanAndDeleteHandler)

	// Only allow certain requests in Development environment only
	if env == "DEV" {
		mux.HandleFunc("/api/areyouawake", server.AwakeHandler)
		mux.HandleFunc("/api/populate", server.PopulateHandler)
		mux.HandleFunc("/api/delete", server.ClearTableHandler(server)) // Why?
	}

	if err := http.ListenAndServe(":4000", mux); err != nil {
		fmt.Print(err)
		log.Fatal(style.Warn.Sprint(err))
	}
}
