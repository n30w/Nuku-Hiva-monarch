package main

// Package for anything that handles data

import (
	"database/sql"
	"log"
	"strconv"
)

type Table struct {
	Name string
}

// A column can either be an int or a string. Ints are for IDs.
type Col interface {
	uint64 | string
}

type id uint64
type author string
type body string
type url string
type subreddit string
type mediaUrl string
type name string

type Post struct {
	Id        id
	Name      name
	URL       url
	Subreddit subreddit
	MediaURL  mediaUrl
}

type Comment struct {
	Id        id
	Author    author
	Body      body
	URL       url
	Subreddit subreddit
}

func RetrieveBy[T Col](db *sql.DB, table *Table, retrieval, want *T) T {
	q := "SELECT " + ToString(retrieval) + " FROM " + table.Name + " WHERE " + ToString(want)

	err := db.QueryRow(q).Scan(&retrieval)
	if err == sql.ErrNoRows {
		log.Fatal("no rows returned")
	} else if err != nil {
		log.Fatal(err)
	}
	return *retrieval
}

func ToString(t any) string {
	switch v := t.(type) {
	case string:
		return v
	case id:
		return strconv.Itoa(int(v))
	default:
		return ""
	}
}
