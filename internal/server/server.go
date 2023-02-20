package server

import (
	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/models"
)

type Server struct {
	RedditPosts, RedditComments, DBPosts, DBComments models.DBTable
	Key                                              credentials.Authenticator
	*models.SQL
}

// New returns a new server object.
func New(key credentials.Authenticator, sql *models.SQL) *Server {
	return &Server{
		RedditPosts:    models.NewTable("posts"),
		RedditComments: models.NewTable("comments"),
		DBPosts:        models.NewTable("posts"),
		DBComments:     models.NewTable("comments"),
		Key:            key,
		SQL:            sql,
	}
}
