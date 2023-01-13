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

	err = s.Retrieve(s.DBPosts, s.DBComments)
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

	w.Write([]byte(Result.Sprintf("Successfully updated Planetscale Database\n")))
	ClearTables(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}

// PopulateHandler handles populating tables requests
func (s *Server) PopulateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Result.Sprintf("Successfully populated Planetscale Database\n")))
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	if err := s.Insert(s.RedditPosts.Name, s.RedditPosts.Rows); err != nil {
		fmt.Println(err)
	}

	if err := s.Insert(s.RedditComments.Name, s.RedditComments.Rows); err != nil {
		fmt.Println(err)
	}
}

// ClearTableHandler handles clearing tables requests
func (s *Server) ClearTableHandler(db RelationalDB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(Result.Sprintf("Cleared all rows from all tables\n")))
		if err := db.Update(s.DBPosts, s.RedditPosts, delete); err != nil {
			fmt.Println(err)
		}

		if err := db.Update(s.DBComments, s.RedditComments, delete); err != nil {
			fmt.Println(err)
		}
	}
}
