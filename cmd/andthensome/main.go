package main

import (
	"database/sql"
	"os"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/server"
	"github.com/n30w/andthensome/internal/style"

	_ "github.com/joho/godotenv/autoload"
)

const (
	PleasePopulateIDs = false
	version           = "1.0.3"
)

var (
	env = os.Getenv("ENVIRONMENT")
	err error
	db  *sql.DB

	// Key objects for authenticating with remote services.
	sqlKey    = &credentials.SQLKey{}
	redditKey = &credentials.RedditKey{}
)

func main() {

	db, err = models.Open("mysql", sqlKey)
	if err != nil {
		panic(style.Warn.Sprint(err))
	}

	err := db.Ping()
	if err != nil {
		panic(style.Warn.Sprint(err))
	}

	dbModel := models.NewSQL(db)

	err = server.New(redditKey, dbModel).Start(4000, env)
	if err != nil {
		panic(err)
	}

}
