package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	RedditPosts, RedditComments, DBPosts, DBComments *Table[Row[id, text]]
	Key                                              *Key
	*PlanetscaleDB
}

// UpdateHandler handles updating SQL database requests
func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	err = s.RetrieveSQL(s.DBPosts, s.DBComments)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		log.Fatal(err)
	}

	err = s.UpdateSQL(s.DBPosts, s.RedditPosts, add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		log.Fatal(err)
	}

	err = s.UpdateSQL(s.DBComments, s.RedditComments, add)
	if err != nil {
		w.Write([]byte(fmt.Sprintf("%s", err)))
		log.Fatal(err)
	}

	w.Write([]byte(Result.Sprint("Successfully updated Planetscale Database")))
	ClearTable(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}

// PopulateHandler handles populating tables requests
func (s *Server) PopulateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Result.Sprint("Successfully populated Planetscale Database")))
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	if err := s.insertToSQL(s.RedditPosts.Name, s.RedditPosts.Rows[:]); err != nil {
		fmt.Println(err)
	}

	if err := s.insertToSQL(s.RedditComments.Name, s.RedditComments.Rows[:]); err != nil {
		fmt.Println(err)
	}
}

// ClearTableHandler handles clearing tables requests
func (s *Server) ClearTableHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte(Result.Sprint("Cleared all rows from all tables")))
	if err := s.UpdateSQL(s.DBPosts, s.RedditPosts, delete); err != nil {
		fmt.Println(err)
	}

	if err := s.UpdateSQL(s.DBComments, s.RedditComments, delete); err != nil {
		fmt.Println(err)
	}
}
