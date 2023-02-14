package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/reddit"
	"github.com/n30w/andthensome/internal/sql"
	"github.com/n30w/andthensome/internal/style"
)

type Server struct {
	RedditPosts, RedditComments, DBPosts, DBComments *models.Table[models.Row[id, text]]
	Key                                              *Key
	*sql.PlanetscaleDB
}

// UpdateHandler handles updating SQL database requests
func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	reddit.GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	err = s.Retrieve(some, s.DBPosts, s.DBComments)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err)))
		log.Fatal(err)
	}

	err = s.Update(s.DBPosts, s.RedditPosts, add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err)))
		log.Fatal(err)
	}

	err = s.Update(s.DBComments, s.RedditComments, add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err)))
		log.Fatal(err)
	}

	w.Write([]byte(style.Result.Sprintf("Successfully updated Planetscale Database\n")))
	ClearTables(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}

// PopulateHandler handles populating tables requests
func (s *Server) PopulateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(style.Result.Sprintf("Successfully populated Planetscale Database\n")))
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	if err := s.Insert(s.RedditPosts.Name, s.RedditPosts.Rows); err != nil {
		fmt.Println(err)
	}

	if err := s.Insert(s.RedditComments.Name, s.RedditComments.Rows); err != nil {
		fmt.Println(err)
	}
}

// AwakeHandler is a route that is used in development.
// Testing uses this route to check if the server is reachable.
func (s *Server) AwakeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(style.Result.Sprintf("Yes, I am awake and accessible. Nice to see you.\n")))
}

// ScanAndDeleteHandler will scan the database for entries that are duplicate and delete them.
func (s *Server) ScanAndDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := s.ScanAndDelete()
	if err != nil {
		w.Write([]byte(fmt.Sprint(err)))
		fmt.Println(err)
	} else {
		w.Write([]byte(style.Result.Sprintf("Scanned and Deleted\n")))
	}

}

// ClearTableHandler handles clearing tables requests
func (s *Server) ClearTableHandler(db models.RelationalDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(style.Result.Sprintf("Cleared all rows from all tables\n")))
		if err := db.Update(s.DBPosts, s.RedditPosts, delete); err != nil {
			fmt.Println(err)
		}

		if err := db.Update(s.DBComments, s.RedditComments, delete); err != nil {
			fmt.Println(err)
		}
	}
}
