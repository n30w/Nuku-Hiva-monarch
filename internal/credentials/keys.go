package credentials

import (
	"os"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

type Authenticator interface {
	Use() any // Use passes a provided key into a function or method parameter.
}

type RedditKey struct{}

func (k *RedditKey) Use() any {
	return k.use()
}

func (k *RedditKey) use() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       os.Getenv("ID"),
		Secret:   os.Getenv("SECRET"),
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}
}

type SQLKey struct{}

func (k *SQLKey) Use() any {
	return k.use()
}

func (k *SQLKey) use() string {
	return os.Getenv(os.Getenv("ENVIRONMENT"))
}
