package models

// Verb is a type of SQL verb.
type Verb uint8

const (
	Add Verb = iota
	Delete
)

// Amount represents an amount of objects.
type Amount uint8

const (
	All Amount = iota
	Some
	Distinct
)
