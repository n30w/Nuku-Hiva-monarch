package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // The underscore on imports autoloads the dependency. Do not need to call something like "godotenv.Load()"
	_ "github.com/joho/godotenv/autoload"
)

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

	for _, v := range table.Rows {
		inserts = append(inserts, insertion)
		params = append(
			params,
			v.Col1,
			v.Col2,
			v.Col3,
			v.Col4,
			v.Col5,
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
