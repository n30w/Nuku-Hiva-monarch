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

	log.Printf("\n%d rows created in %s", rows, tableName)

	return nil
}

// RetrieveSQL stores the most recent 10 rows from a PlanetscaleDB table
// into the parameterized table.
func (p *PlanetscaleDB) RetrieveSQL(table *Table[Row[id, text]]) error {

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

	return nil
}

// UpdateSQL compares the planetscale table and the reddit table,
// and updates the planetscale database accordingly.
func (p *PlanetscaleDB) UpdateSQL(planetscale, reddit *Table[Row[id, text]]) {

	if planetscale.Name != reddit.Name {
		log.Fatal("these databases are not the same")
	}

	entriesToAdd := 0
	for i := 0; i < 10; i++ {
		if planetscale.Rows[0].Col3 == reddit.Rows[i].Col3 {
			entriesToAdd = i
		}
	}

	if entriesToAdd == 0 {
		log.Println("No new rows must be added to " + planetscale.Name)
	} else {
		lastID := p.GetLastID(planetscale.Name)
		var nextID id
		for i := 0; i < entriesToAdd; i++ {
			nextID = lastID + id(i+1)
			reddit.Rows[i].Col1 = nextID
		}
		p.InsertToSQL(planetscale.Name, reddit.Rows[0:entriesToAdd])
	}
}

// GetLastID makes a query to the SQL database and returns the most latest ID
func (p *PlanetscaleDB) GetLastID(name string) id {
	rows, err := p.Query("SELECT MAX(id) FROM " + name)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var max int
	for rows.Next() {
		err := rows.Scan(&max)
		if err != nil {
			log.Fatal(err)
		}
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	return id(max)
}
