package main

type Table string

type Col interface {
	string
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
	Name      name
	URL       url
	Subreddit subreddit
	Author    author
	Body      body
}
