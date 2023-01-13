package main

import (
	"strconv"
	"testing"
)

func Test_entriesToAdd(t *testing.T) {
	planetscale := [100]*Row[id, text]{}
	reddit := [100]*Row[id, text]{}

	planetscale[0] = &Row[id, text]{Col3: "18"}

	t.Run("addition of n entries", func(t *testing.T) {
		want := 5
		reddit[want] = &Row[id, text]{Col3: "18"}

		// Populate with arbitrary data
		for i := 0; i < want; i++ {
			reddit[i] = &Row[id, text]{Col3: text(strconv.Itoa(i))}
		}

		got := entriesToAdd(
			planetscale[:ResultsPerRedditRequest],
			reddit[:ResultsPerRedditRequest],
		)

		if got != want {
			t.Errorf("got %d entries to add, want %d", got, want)
		}
	})

	t.Run("addition of 0 entries", func(t *testing.T) {
		want := 0
		(*reddit[0]).Col3 = "18"
		got := entriesToAdd(
			planetscale[:ResultsPerRedditRequest],
			reddit[:ResultsPerRedditRequest],
		)

		if got != want {
			t.Errorf("got %d entries, want %d", got, want)
		}
	})

	t.Run("check nil exit", func(t *testing.T) {
		want := 1
		(*reddit[0]).Col3 = "1"
		got := entriesToAdd(
			planetscale[:ResultsPerRedditRequest],
			reddit[4:ResultsPerRedditRequest],
		)

		if got != want {
			t.Errorf("got %d entries, want only %d", got, want)
		}
	})

}
