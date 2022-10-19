package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // The underscore on imports autoloads the dependency. Do not need to call something like "godotenv.Load()"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Key represents credentials used to login to APIs
type Key struct{}

// NewKey returns a new key given environment variables
func (k *Key) NewKey() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
}

type DBModel interface {
	UpdateSQL(table *Table[Row[id, text]]) error
	InsertToSQL(table *Table[Row[id, text]]) error
}

type PlanetscaleDB struct {
	*sql.DB
}

// UpdateSQL updates table rows for any changes
func (p *PlanetscaleDB) UpdateSQL(table *Table[Row[id, text]]) error {
	return nil
}

// InsertToSql uploads table data to SQL
func (p *PlanetscaleDB) InsertToSQL(table *Table[Row[id, text]]) error {

	var query string
	var inserts []string
	var params []interface{}
	insertion := "(?, ?, ?, ?, ?)"

	switch table.Name {
	case "posts":
		query = "INSERT INTO posts (id, name, url, subreddit, media_url) VALUES "
	case "comments":
		query = "INSERT INTO comments (id, author, body, url, subreddit) VALUES "
	default:
		return errors.New("not a valid table name")
	}

	for _, row := range table.Rows {
		if row == nil {
			break
		}
		inserts = append(inserts, insertion)
		params = append(
			params,
			row.Col1,
			row.Col2,
			row.Col3,
			row.Col4,
			row.Col5,
		)
	}

	qv := strings.Join(inserts, ",")

	query += qv

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancelFunc()

	statement, err := p.PrepareContext(ctx, query)

	if err != nil {
		log.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	defer statement.Close()

	res, err := statement.ExecContext(ctx, params...)
	if err != nil {
		log.Printf("Error %s when inserting row into products table", err)
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}

	log.Printf("\n%d rows created ", rows)

	return nil
}
