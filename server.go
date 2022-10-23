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
	w.Write([]byte("Successfully updated Planetscale Database"))

	GrabSaved(s.RedditPosts, s.RedditComments, s.Key, 1)

	err := s.Psdb.RetrieveSQL(s.DBPosts)
	if err != nil {
		log.Fatalln(err)
	}
	err = s.Psdb.RetrieveSQL(s.DBComments)
	if err != nil {
		log.Fatalln(err)
	}

	s.Psdb.UpdateSQL(s.DBPosts, s.RedditPosts)
	s.Psdb.UpdateSQL(s.DBComments, s.RedditComments)

}
