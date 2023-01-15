package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // The underscore will autoload the dependency.
	// Do not need to call something like "godotenv.Load()"
	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// verb is a type of SQL verb
type verb uint8

const (
	add verb = iota
	delete
)

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

type PlanetscaleDB struct {
	*sql.DB
}

// Insert creates a string consisting of all the rows
// in a given table, and executes the query, inserting
// the items into the Planetscale database.
func (p *PlanetscaleDB) Insert(tableName string, tableRows Rows) error {
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
		return errors.New("table does not exist")
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

// Delete deletes all rows from a specified table.
func (p *PlanetscaleDB) Delete(tableName string) error {
	query, err := p.Query("DELETE FROM " + tableName)
	if err != nil {
		return err
	}
	defer query.Close()
	return nil
}

// Retrieve stores the most recent n rows from a PlanetscaleDB table
// into the parameterized table.
func (p *PlanetscaleDB) Retrieve(tables ...DBTable) error {
	for _, table := range tables {
		rows, err := p.Query("SELECT * FROM " + table.Name + " ORDER BY id DESC LIMIT " + strconv.Itoa(ResultsPerRedditRequest))
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

// Update compares the planetscale table and the reddit table,
// and updates the planetscale database accordingly. This is essentially
// a sync function that synchronizes the planetscale database
// and the Reddit saved posts list
func (p *PlanetscaleDB) Update(planetscale, reddit DBTable, v verb) error {

	if planetscale.Name != reddit.Name {
		return errors.New(Warn.Sprintf("these tables are not the same"))
	}

	// TODO make a test for case: add
	switch v { // TODO go routine optimization can occur here
	case add:
		msg := Information.Sprint("No new rows must be added to " + planetscale.Name)
		mostRecentIDOnPlanetscale := p.getLastId(planetscale.Name)
		entries := entriesToAdd(
			planetscale.Rows[0:ResultsPerRedditRequest],
			reddit.Rows[0:ResultsPerRedditRequest],
		)

		if entries == 0 {
			log.Print(msg)
			return nil
		} else {
			for i := 0; i < entries; i++ {
				reddit.Rows[i].Col1 = mostRecentIDOnPlanetscale + id(entries-i)
			}
			p.Insert(planetscale.Name, reddit.Rows)
		}
	case delete:
		msg := Information.Sprint("Deleted rows from SQL tables")
		if err := p.Delete(planetscale.Name); err != nil {
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
		if planetscale[0].Col3 == reddit[i].Col3 || reddit[i+1] == nil {
			return i
		}
	}
	return 0
}
