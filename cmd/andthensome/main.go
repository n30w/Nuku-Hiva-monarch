package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
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
	key := credentials.NewKey()

	var err error
	db, err = sql.Open("mysql", key.SQLKey())
	if err != nil {
		panic(style.Warn.Sprint(err))
	}

	if err := db.Ping(); err != nil {
		panic(style.Warn.Sprint(err))
	}

	server := &Server{
		RedditPosts:    models.NewTable("posts"),
		RedditComments: models.NewTable("comments"),
		DBPosts:        models.NewTable("posts"),
		DBComments:     models.NewTable("comments"),
		Key:            key,
		PlanetscaleDB:  &PlanetscaleDB{db},
	}

	log.Print(style.Start.Sprintf("Starting andthensome %s %s", version, env))
	log.Print(style.Start.Sprint("Server listening on :4000"))

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
