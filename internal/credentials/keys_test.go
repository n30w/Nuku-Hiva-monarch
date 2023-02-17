package credentials

import (
	"reflect"
	"testing"

	"github.com/vartanbeno/go-reddit/v2/reddit"
)

func requireAuth(any any) any {
	return any
}

// manipulateDataMock mocks a function that manipulates or transforms data, like
// an SQL insert function or a function that sends credential information over
// to Reddit.
func manipulateDataMock(auth Authenticator) any {
	return requireAuth(auth.Use())
}

type redditKeyMock struct{}

func (rkm *redditKeyMock) Use() any {
	return rkm.use()
}

func (rkm *redditKeyMock) use() *reddit.Credentials {
	return &reddit.Credentials{
		ID:       "id",
		Secret:   "secret",
		Username: "username",
		Password: "password",
	}
}

type sqlKeyMock struct{}

func (skm *sqlKeyMock) Use() any {
	return skm.use()
}

func (skm *sqlKeyMock) use() string {
	return "sql"
}

func TestAuthenticatorInterface(t *testing.T) {

	t.Run("reddit key Use", func(t *testing.T) {

		rk := &redditKeyMock{}

		got := manipulateDataMock(rk)
		want := &reddit.Credentials{}

		if reflect.TypeOf(got) != reflect.TypeOf(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})

	t.Run("sql key Use", func(t *testing.T) {

		sk := &sqlKeyMock{}

		got := manipulateDataMock(sk)
		var want string

		if reflect.TypeOf(got) != reflect.TypeOf(want) {
			t.Errorf("got %s, want %s", got, want)
		}
	})
}
