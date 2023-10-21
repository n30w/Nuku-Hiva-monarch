package main

import (
	"database/sql"
	"flag"
	"os"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/server"

	_ "github.com/joho/godotenv/autoload"
)

var (
	env = os.Getenv("ENVIRONMENT")
	err error

	rk = &credentials.RedditKey{}
	sk = &credentials.SQLKey{}
	db = models.NewSQL(&sql.DB{})
)

func main() {
	s := server.New(rk, sk, db).Initialize("mysql")

	oneShotMode := *flag.Bool("osm", false, "launch in oneshot mode")

	flag.Parse()

	switch oneShotMode {
	case true:
		err = s.Start(4000, env)
	case false:
		err = s.OneShot()
	}

	if err != nil {
		panic(err)
	}

}
