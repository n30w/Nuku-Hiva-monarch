package main

import (
	"database/sql"
	"os"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

var (
	env = os.Getenv("ENVIRONMENT")
	err error

	db = models.NewSQL(&sql.DB{})
	rk = &credentials.RedditKey{}
	sk = &credentials.SQLKey{}
)

func main() {
	err = server.
		New(rk, sk, db).
		Initialize("mysql").
		Start(4000, env)

	if err != nil {
		panic(err)
	}

}
