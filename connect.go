package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // The underscore will autoload the dependency.
	// Do not need to call something like "godotenv.Load()"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// verb is a type of SQL verb
type verb string

// Key represents credentials used to log in to APIs
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

type Rows []*Row[id, text]
type DBTable *Table[Row[id, text]]

type DBModel interface {
	insertToSQL(tableName string, tableRows Rows) error
	RetrieveSQL(tables ...DBTable) error
	UpdateSQL(planetscale, reddit DBTable, v verb) error
}

type PlanetscaleDB struct {
	*sql.DB
}

// InsertToSQL creates a string consisting of all the rows
// in a given table, and executes the query, inserting
// the items into the Planetscale database.
func (p *PlanetscaleDB) insertToSQL(tableName string, tableRows Rows) error {
	var query string
	var inserts []string
	var params []interface{}
	insertion := "(?, ?, ?, ?, ?)"

	switch tableName {
	case "posts":
		query = "INSERT INTO posts (id, name, url, subreddit, media_url) VALUES "
	case "comments":
		query = "INSERT INTO comments (id, author, body, url, subreddit) VALUES "
	default:
		return errors.New("not a valid table name")
	}

	for _, row := range tableRows {
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
		log.Print(Warn.Sprintf("Error %s when preparing SQL statement", err))
		return err
	}

	defer statement.Close()

	res, err := statement.ExecContext(ctx, params...)
	if err != nil {
		log.Print(Warn.Sprintf("Error %s when inserting row into products table", err))
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		log.Print(Warn.Sprintf("Error %s when finding rows affected", err))
		return err
	}

	log.Print(Result.Sprintf("%d rows created in %s", rows, tableName))

	return nil
}

func (p *PlanetscaleDB) DeleteRowsFromSQL(tableName string) error {
	query, err := p.Query("DELETE FROM " + tableName)
	if err != nil {
		return err
	}
	defer query.Close()
	return nil
}

// RetrieveSQL stores the most recent 10 rows from a PlanetscaleDB table
// into the parameterized table.
func (p *PlanetscaleDB) RetrieveSQL(tables ...DBTable) error {
	for _, table := range tables {
		rows, err := p.Query("SELECT * FROM " + table.Name + " ORDER BY id DESC LIMIT 10")
		if err != nil {
			return err
		}
		defer rows.Close()

		i := 0
		for rows.Next() {
			table.Rows[i] = &Row[id, text]{}
			err := rows.Scan(
				&table.Rows[i].Col1,
				&table.Rows[i].Col2,
				&table.Rows[i].Col3,
				&table.Rows[i].Col4,
				&table.Rows[i].Col5,
			)
			if err != nil {
				return err
			}
			i++
		}

		if err = rows.Err(); err != nil {
			return err
		}
	}
	return nil
}

// UpdateSQL compares the planetscale table and the reddit table,
// and updates the planetscale database accordingly. This is essentially
// a sync function that synchronizes the planetscale database
// and the Reddit saved posts list
func (p *PlanetscaleDB) UpdateSQL(planetscale, reddit DBTable, v verb) error {

	if planetscale.Name != reddit.Name {
		return errors.New(Warn.Sprintf("these tables are not the same"))
	}

	switch v {
	case "ADD":
		msg := Information.Sprint("No new rows must be added to " + planetscale.Name)
		mostRecentIDOnPlanetscale := p.getLastId(planetscale.Name)
		entries := entriesToAdd(
			planetscale.Rows[0:ResultsPerRedditRequest+1],
			reddit.Rows[0:ResultsPerRedditRequest+1],
		)

		if entries == 0 {
			log.Print(msg)
			return nil
		} else {
			for i := 0; i < entries; i++ {
				reddit.Rows[i].Col1 = mostRecentIDOnPlanetscale + id(entries-i)
			}
			p.insertToSQL(planetscale.Name, reddit.Rows[0:entries])
		}
	case "DELETE":
		msg := Information.Sprint("Deleted rows from SQL tables")
		if err := p.DeleteRowsFromSQL(planetscale.Name); err != nil {
			return errors.New(Warn.Sprintf("Could not delete tables: %s", err))
		} else {
			log.Print(msg)
		}

	default:
		return errors.New(Warn.Sprint("no operation provided in UpdateSQL()"))
	}

	return nil
}

// getLastID makes a query to the SQL database and returns the latest ID
func (p *PlanetscaleDB) getLastId(name string) id {
	rows, err := p.Query("SELECT MAX(id) FROM " + name)
	if err != nil {
		log.Fatal(Warn.Sprint(err))
	}
	defer rows.Close()

	var max int
	for rows.Next() {
		err := rows.Scan(&max)
		if err != nil {
			log.Fatal(Warn.Sprint(err))
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(Warn.Sprint(err))
	}

	return id(max)
}

// entriesToAdd compares two rows, one from Planetscale and one from Reddit.
// It returns an integer, which represents the number of rows to update.
func entriesToAdd(planetscale, reddit []*Row[id, text]) int {
	for i := 0; i < ResultsPerRedditRequest; i++ {
		if planetscale[0].Col3 == reddit[i].Col3 {
			return i
		}
	}
	return 0
}

// compareIncrement compares two tables and returns a slice of
// ids at which to delete in order ot update the planetscale database
// func (p *PlanetscaleDB) compareIndex(planetscale, reddit *Table[Row[id, text]]) []id {
// 	idsToUpdate := make([]id, 0)

// 	return idsToUpdate
// }
