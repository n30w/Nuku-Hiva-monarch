package main

import (
	"fmt"
	"log"
	"net/http"
)

type Server struct {
	RedditPosts, RedditComments, DBPosts, DBComments *Table[Row[id, text]]
	Key                                              *Key
	Psdb                                             *PlanetscaleDB
	Environment                                      string
}

func (s *Server) UpdateHandler(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Write([]byte(Result.Sprint("Successfully updated Planetscale Database")))
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	var add verb = "ADD"
	err = s.Psdb.RetrieveSQL(s.DBPosts, s.DBComments)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	err = s.Psdb.UpdateSQL(s.DBPosts, s.RedditPosts, add)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	err = s.Psdb.UpdateSQL(s.DBComments, s.RedditComments, add)
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

	ClearTable(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}

func (s *Server) PopulateHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(Result.Sprint("Successfully populated Planetscale Database")))
	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	if err := s.Psdb.insertToSQL(s.RedditPosts.Name, s.RedditPosts.Rows[:]); err != nil {
		fmt.Println(err)
	}
	if err := s.Psdb.insertToSQL(s.RedditComments.Name, s.RedditComments.Rows[:]); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) ClearTableHandler(w http.ResponseWriter, r *http.Request) {
	var delete verb = "DELETE"
	w.Write([]byte(Result.Sprint("Cleared all rows from all tables")))
	if err := s.Psdb.UpdateSQL(s.DBPosts, s.RedditPosts, delete); err != nil {
		fmt.Println(err)
	}
	if err := s.Psdb.UpdateSQL(s.DBComments, s.RedditComments, delete); err != nil {
		fmt.Println(err)
	}
}
