package models

// Verb is a type of SQL verb.
type Verb verb

type verb uint8

const (
	Add Verb = iota
	Delete
)

// Amount represents an amount of objects that are queried to be retrieved.
type Amount amount

type amount uint8

const (
	All Amount = iota
	Some
	Distinct
)

type RelationalDB interface {
	Insert(tableName string, tableRows Rows) error
	Delete(tableName string) error
	Retrieve(amount Amount, tables ...DBTable) error
	Update(planetscale, reddit DBTable, v Verb) error
}
