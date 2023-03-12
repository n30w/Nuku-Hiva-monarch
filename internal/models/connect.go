package models

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // The underscore will autoload the dependency.
	// Do not need to call something like "godotenv.Load()"
	"github.com/n30w/andthensome/internal/credentials"
	"github.com/n30w/andthensome/internal/style"
)

const ResultsPerRedditRequest = 50

// SQL wraps an SQL Database tools provided by Go standard pkg,
// since the functionality of the remote database I'm using, which
// is Planetscale, is exactly the same.
type SQL struct {
	*sql.DB
}

// NewSQL creates and returns a new SQL object.
func NewSQL(db *sql.DB) *SQL {
	return &SQL{DB: db}
}

// Open opens a new connection to a database with an SQL Driver.
func Open(driverName string, auth credentials.Authenticator) (*sql.DB, error) {
	db, err := sql.Open("mysql", auth.Use().(string))
	return db, err
}

type RelationalDB interface {
	Insert(tableName string, tableRows Rows) error
	Delete(tableName string) error
	Retrieve(amount Amount, tables ...DBTable) error
	Update(planetscale, reddit DBTable, v Verb) error
}

// Operator performs CRUD operations.
type Operator interface {
	Create() error
	Read() error
	Update() error
	Delete() error
}

// Insert creates a string consisting of all the rows
// in a given table, and executes the query, inserting
// the items into the Planetscale database.
func (p *SQL) Insert(tableName string, tableRows Rows) error {
	var query string
	var inserts []string
	var params []interface{}
	insertion := "(?, ?, ?, ?)"

	switch tableName {
	case "posts":
		query = "INSERT INTO posts (name, url, subreddit, media_url) VALUES "
	case "comments":
		query = "INSERT INTO comments (author, body, url, subreddit) VALUES "
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
		log.Print(style.Warn.Sprintf("Error %s when preparing SQL statement", err))
		return err
	}

	defer statement.Close()

	res, err := statement.ExecContext(ctx, params...)
	if err != nil {
		log.Print(style.Warn.Sprintf("Error %s when inserting row into products table", err))
		return err
	}

	rows, err := res.RowsAffected()

	if err != nil {
		log.Print(style.Warn.Sprintf("Error %s when finding rows affected", err))
		return err
	}

	log.Print(style.Result.Sprintf("%d rows created in %s", rows, tableName))

	return nil
}

// Delete deletes all rows from a specified table.
func (p *SQL) Delete(tableName string) error {
	query, err := p.Query("DELETE FROM " + tableName)
	if err != nil {
		return err
	}
	defer query.Close()
	return nil
}

// Retrieve stores the most recent n rows from an SQL table
// into the parameterized table.
func (p *SQL) Retrieve(amount Amount, tables ...DBTable) error {
	for _, table := range tables {

		var rows *sql.Rows
		var err error

		switch amount {
		case All:
			rows, err = p.Query("SELECT * FROM " + table.Name + " ORDER BY id DESC LIMIT 10000")
		case Some:
			rows, err = p.Query("SELECT * FROM " + table.Name + " ORDER BY id DESC LIMIT " + strconv.Itoa(ResultsPerRedditRequest))
		case Distinct:
			switch table.Name {
			case "posts": // select distinct ... from ... group by ...
				rows, err = p.Query("SELECT id, name, url, subreddit, media_url FROM (SELECT name, url, subreddit, media_url, MAX(id) id FROM `posts` GROUP BY name, url, subreddit, media_url) A ORDER BY id")
			case "comments":
				rows, err = p.Query("SELECT id, author, body, url, subreddit FROM (SELECT author, body, url, subreddit, MAX(id) id FROM `comments` GROUP BY author, body, url, subreddit) A ORDER BY id")

			default:
				return errors.New("invalid table name")
			}

		default:
			return errors.New("connect.go: invalid amount parameter")
		}

		if err != nil {
			return err
		}

		defer rows.Close()

		i := 0

		for rows.Next() {
			table.Rows[i] = &Row[id, text]{}
			err = rows.Scan(
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
func (p *SQL) Update(planetscale, reddit DBTable, verb Verb) error {

	if planetscale.Name != reddit.Name {
		return errors.New(style.Warn.Sprintf("these tables are not the same"))
	}

	// TODO make a test for case: add
	switch verb { // TODO go routine optimization can occur here
	case Add:
		msg := style.Information.Sprint("No new rows must be added to " + planetscale.Name)

		// insertion contains the new rows to insert into the database.
		insertion := Table[Row[id, text]]{Name: planetscale.Name}

		// inventory is a map of current rows on the SQL database.
		inventory := make(map[text]bool)

		if err := p.Retrieve(Some, planetscale); err != nil {
			return err
		}

		// This loop adds the planetscale rows to the inventory to check later.
		for _, row := range planetscale.Rows {
			if row == nil {
				continue
			}

			inventory[row.Col3] = true
		}

		// index keeps track of current row.
		index := 0
		for _, row := range reddit.Rows {
			if row == nil {
				continue
			}
			// if the row doesn't exist in inventory, add it to the insertion table.
			if !inventory[row.Col3] {
				insertion.Rows[index] = &Row[id, text]{
					// Col1: row.Col1, // Might need to change the ID.
					row.Col1,
					row.Col2,
					row.Col3,
					row.Col4,
					row.Col5,
				}

				index++
			}
		}

		// if there are no new rows, tell the user.
		if index == 0 {
			log.Print(msg)
			return nil
		}

		// Finally, insert new rows into SQL database.
		if err := p.Insert(planetscale.Name, insertion.Rows); err != nil {
			return err
		}

	case Delete:
		msg := style.Information.Sprint("Deleted rows from SQL tables")
		if err := p.Delete(planetscale.Name); err != nil {
			return errors.New(style.Warn.Sprintf("Could not delete tables: %s", err))
		} else {
			log.Print(msg)
		}

	default:
		return errors.New(style.Warn.Sprint("no operation provided in UpdateSQL()"))
	}

	return nil
}

// ScanAndDelete retrieves entries from the SQL database,
// and deletes the duplicate ones, regardless of ID number.
func (p *SQL) ScanAndDelete() error {

	if err := p.deleteDuplicates(); err != nil {
		return err
	}

	return nil
}

// deleteDuplicates deletes duplicate entries in the database
func (p *SQL) deleteDuplicates() error {

	// Checkout this thread:
	// https://dba.stackexchange.com/questions/19511/getting-unique-names-when-the-ids-are-different-distinct

	// Query to delete duplicate rows
	// DELETE FROM comments
	// WHERE id NOT IN (
	// SELECT id FROM (SELECT MAX(id) id FROM `comments` GROUP BY body) A ORDER BY id
	// );

	if _, err := p.Query("DELETE FROM posts WHERE id NOT IN (SELECT id FROM (SELECT MAX(id) id FROM `posts` GROUP BY url) A ORDER BY id)"); err != nil {
		return err
	}

	log.Println("Deleted duplicate rows from posts")

	if _, err := p.Query("DELETE FROM posts WHERE url LIKE ''"); err != nil {
		return err
	}

	log.Println("Deleted rows with empty content from posts")

	if _, err := p.Query("DELETE FROM comments WHERE id NOT IN (SELECT id FROM (SELECT MAX(id) id FROM `comments` GROUP BY body) A ORDER BY id)"); err != nil {
		return err
	}

	log.Println("Deleted duplicate rows from comments")

	if _, err := p.Query("DELETE FROM comments WHERE url LIKE ''"); err != nil {
		return err
	}

	log.Println("Deleted rows with empty content from comments")

	return nil
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
