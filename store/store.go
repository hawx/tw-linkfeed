// Package store implements a data-structure for storing lists of items which
// expire after a certain time.
package store

import (
	"github.com/hawx/tw-linkfeed/stream"

	"container/ring"
	"sync"
	"time"
)

type Store interface {
	Add(tweet stream.Tweet)
	Latest() []stream.Tweet
}

type store struct {
	swap   time.Time
	bucket *ring.Ring
	mutex  *sync.RWMutex
}

type bucket []stream.Tweet

func New(n int, interval time.Duration) Store {
	return &store{
		swap:   time.Now().Add(interval),
		bucket: ring.New(n),
		mutex:  &sync.RWMutex{},
	}
}

func (s store) Add(tweet stream.Tweet) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if time.Now().After(s.swap) {
		s.bucket = s.bucket.Next()
		s.bucket.Value = nil
	}

	if s.bucket.Value == nil {
		s.bucket.Value = bucket{tweet}
		return
	}

	s.bucket.Value = append(s.bucket.Value.(bucket), tweet)
}

func (s store) Latest() []stream.Tweet {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	tweets := []stream.Tweet{}
	s.bucket.Do(func(value interface{}) {
		if value != nil {
			tweets = append(tweets, value.(bucket)...)
		}
	})

	return tweets
}
