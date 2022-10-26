package main

import (
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

	var operation verb = "ADD"
	err = s.Psdb.RetrieveSQL(s.DBPosts, s.DBComments)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Psdb.UpdateSQL(s.DBPosts, s.RedditPosts, operation)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Psdb.UpdateSQL(s.DBComments, s.RedditComments, operation)
	if err != nil {
		log.Fatal(err)
	}

	ClearTable(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}
