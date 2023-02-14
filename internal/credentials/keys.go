package credentials

import (
	"os"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

// Key represents credentials used to log in to APIs.
type Key struct{}

// NewKey creates a new key object that can access keys.
func NewKey() *Key {
	return &Key{}
}

// RedditKey returns reddit credentials based on environment variables.
func (k *Key) RedditKey() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
}

// SQLKey returns an SQL access key via environment variables.
func (k *Key) SQLKey() string {
	return os.Getenv(os.Getenv("ENVIRONMENT"))
}
