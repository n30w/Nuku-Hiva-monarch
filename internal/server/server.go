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

// New returns a new server prototype.
func New(redditKey, dbKey credentials.Authenticator, sqlModel *models.SQL) *Server {
	return &Server{
		RedditPosts:    models.NewTable("posts"),
		RedditComments: models.NewTable("comments"),
		DBPosts:        models.NewTable("posts"),
		DBComments:     models.NewTable("comments"),
		RedditKey:      redditKey,
		DBKey:          dbKey,
		Sql:            sqlModel,
	}
}

type Server struct {
	// Table models used for data manipulation and processing.
	RedditPosts, RedditComments, DBPosts, DBComments models.DBTable

	// Credentials like secrets stored in environment variables that
	// are functionally retrievable by these objects.
	RedditKey, DBKey credentials.Authenticator

	// Sql object to perform sql operations on a remote database.
	Sql *models.SQL
}

// update retrieves the saved reddit posts and comments, and updates the SQL database
// according to whether or not the database and the newly retrieved objects match.
// Finally, it clears tables in order for later use.
func (s *Server) update() error {
	var err error

	err = reddit.Saved(s.RedditPosts, s.RedditComments, s.RedditKey)
	if err != nil {
		return err
	}

	err = s.Sql.Retrieve(models.Some, s.DBPosts, s.DBComments)
	if err != nil {
		return err
	}

	err = s.Sql.Update(s.DBPosts, s.RedditPosts, models.Add)
	if err != nil {
		return err
	}

	err = s.Sql.Update(s.DBComments, s.RedditComments, models.Add)
	if err != nil {
		return err
	}

	models.ClearTables(s.RedditPosts, s.RedditComments, s.DBPosts, s.DBComments)

	return nil
}

// scanDelete will scan the database for entries that are duplicate and delete them.
func (s *Server) scanDelete() error {
	if err := s.Sql.ScanAndDelete(); err != nil {
		return err
	}
	return nil
}

// populate grabs all the saved content from Reddit and adds it to the database.
func (s *Server) populate() error {
	reddit.Saved(s.RedditPosts, s.RedditComments, s.RedditKey)

	if err := s.Sql.Insert(s.RedditPosts.Name, s.RedditPosts.Rows); err != nil {
		return err
	}

	if err := s.Sql.Insert(s.RedditComments.Name, s.RedditComments.Rows); err != nil {
		return err
	}

	return nil
}

// Start establishes server routes using handlers and starts the server.
func (s *Server) Start(port int, env string) error {

	log.Print(style.Start.Sprintf("Starting andthensome on %s", env))
	mux := http.NewServeMux()

	mux.HandleFunc("/", homeHandler)

	mux.HandleFunc("/api/update", updateHandler(s.update))
	mux.HandleFunc("/api/scanDelete", scanDeleteHandler(s.scanDelete))

	// Only allow certain requests in Development environment only
	if env == "DEV" {
		mux.HandleFunc("/api/populate", populateHandler(s.populate))
	}

	log.Print(style.Start.Sprintf("Server listening on port %d", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), mux); err != nil {
		log.Fatal(style.Warn.Sprint(err))
	}

	return nil
}

// Initialize initializes a connection to a database using the server's sql database object.
func (s *Server) Initialize(driverName string) *Server {
	var err error

	s.Sql.DB, err = models.Open(driverName, s.DBKey)
	if err != nil {
		panic(style.Warn.Sprint(err))
	}

	err = s.Sql.DB.Ping()
	if err != nil {
		panic(style.Warn.Sprint(err))
	}

	return s
}
