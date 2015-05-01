package store

import (
	"testing"
	"time"

	"hawx.me/code/tw-linkfeed/stream"
)

func TestStore(t *testing.T) {
	store := New(3, time.Second)

	for i := 0; i < 8; i++ {
		store.Add(stream.Tweet{Id: int64(i)})
		time.Sleep(time.Second / 2)
	}

	all := store.Latest()

	if len(all) != 4 {
		t.Log(all)
		t.FailNow()
	}

	for i, tweet := range all {
		if tweet.Id != int64(7-i) {
			t.Log(all)
			t.FailNow()
		}
	}
}
