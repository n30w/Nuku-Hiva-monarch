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
	w.Write([]byte(Result.Sprint("Successfully updated Planetscale Database")))

	GrabSaved(s.RedditPosts, s.RedditComments, s.Key)

	if PleasePopulateIDs {
		s.Psdb.InsertToSQL(s.RedditPosts.Name, s.RedditPosts.Rows[:])
	} else {
		if err := s.Psdb.RetrieveSQL(s.DBPosts, s.DBComments); err != nil {
			log.Fatal(Warn.Sprint(err))
		}

		s.Psdb.UpdateSQL(s.DBPosts, s.RedditPosts)
		s.Psdb.UpdateSQL(s.DBComments, s.RedditComments)
	}

	ClearTable(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)
}
