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
	InsertToSQL(table *Table[Row[id, text]]) error
	RetrieveSQL(table *Table[Row[id, text]]) error
	UpdateSQL(table *Table[Row[id, text]]) error
}

type PlanetscaleDB struct {
	*sql.DB
}

// InsertToSQL uploads table data to SQL
func (p *PlanetscaleDB) InsertToSQL(tableName string, tableRows []*Row[id, text]) error {
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

// RetrieveSQL stores the most recent 10 rows from a PlanetscaleDB table
// into the parameterized table.
func (p *PlanetscaleDB) RetrieveSQL(tables ...*Table[Row[id, text]]) error {
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
// and updates the planetscale database accordingly.
func (p *PlanetscaleDB) UpdateSQL(planetscale, reddit *Table[Row[id, text]]) {

	msg := Information.Sprint("No new rows must be added to " + planetscale.Name)

	if planetscale.Name != reddit.Name {
		log.Fatal(Warn.Sprintf("these tables are not the same"))
	}

	mostRecentIDOnPlanetscale := p.getLastID(planetscale.Name)

	entriesToAdd := p.compareRows(planetscale, reddit)

	if entriesToAdd == 0 {
		log.Print(msg)
		return
	} else {
		for i := 0; i < entriesToAdd; i++ {
			reddit.Rows[i].Col1 = mostRecentIDOnPlanetscale + id(entriesToAdd-i)
		}
		p.InsertToSQL(planetscale.Name, reddit.Rows[0:entriesToAdd])
	}
}

// compareRows compares two rows, one from Planetscale and one from Reddit.
// It returns an integer, which represents the number of rows to update.
func (p *PlanetscaleDB) compareRows(planetscale, reddit *Table[Row[id, text]]) int {
	for i := 0; i < ResultsPerRedditRequest; i++ {
		if planetscale.Rows[0].Col3 == reddit.Rows[i].Col3 {
			return i
		}
	}
	return 0
}

// GetLastID makes a query to the SQL database and returns the most latest ID
func (p *PlanetscaleDB) getLastID(name string) id {
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
