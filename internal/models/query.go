package models

type verb uint8
type amount uint8

// Verb is a type of SQL verb.
type Verb verb

const (
	Add Verb = iota
	Delete
)

// Amount represents an amount of objects that are queried to be retrieved.
type Amount amount

const (
	All Amount = iota
	Some
	Distinct
)
