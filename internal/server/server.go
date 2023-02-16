package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
	"github.com/n30w/andthensome/internal/reddit"
	"github.com/n30w/andthensome/internal/style"
)

type Server struct {
	RedditPosts, RedditComments, DBPosts, DBComments models.DBTable
	Key                                              *credentials.Key
	*models.SQL
}

// NewServer returns a new server object.
func New(key *credentials.Key, sql *models.SQL) *Server {
	return &Server{
		RedditPosts:    models.NewTable("posts"),
		RedditComments: models.NewTable("comments"),
		DBPosts:        models.NewTable("posts"),
		DBComments:     models.NewTable("comments"),
		Key:            key,
		SQL:            sql,
	}
}

// UpdateHandler handles updating SQL database requests
func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	reddit.GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	err = s.Retrieve(models.Some, s.DBPosts, s.DBComments)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err))) //nolint
		log.Fatal(err)
	}

	err = s.Update(s.DBPosts, s.RedditPosts, models.Add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err))) //nolint
		log.Fatal(err)
	}

	err = s.Update(s.DBComments, s.RedditComments, models.Add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s\n", err))) //nolint
		log.Fatal(err)
	}

	w.Write([]byte(style.Result.Sprintf("Successfully updated database\n"))) //nolint
	models.ClearTables(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}

// PopulateHandler handles populating tables requests
func (s *Server) PopulateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(style.Result.Sprintf("Successfully populated Planetscale Database\n"))) //nolint
	reddit.GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

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
	w.Write([]byte(style.Result.Sprintf("Yes, I am awake and accessible. Nice to see you.\n"))) //nolint
}

// ScanAndDeleteHandler will scan the database for entries that are duplicate and delete them.
func (s *Server) ScanAndDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := s.ScanAndDelete()
	if err != nil {
		w.Write([]byte(fmt.Sprint(err))) //nolint
		fmt.Println(err)
	} else {
		w.Write([]byte(style.Result.Sprintf("Scanned and Deleted\n"))) //nolint
	}

}

// ClearTableHandler handles clearing tables requests
func (s *Server) ClearTableHandler(db models.RelationalDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(style.Result.Sprintf("Cleared all rows from all tables\n"))) //nolint
		if err := db.Update(s.DBPosts, s.RedditPosts, models.Delete); err != nil {
			fmt.Println(err)
		}

		if err := db.Update(s.DBComments, s.RedditComments, models.Delete); err != nil {
			fmt.Println(err)
		}
	}
}
