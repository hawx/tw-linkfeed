// Package store implements a data-structure for storing lists of items which
// expire after a certain time.
package store

import (
	"github.com/hawx/tw-stream"

	"sync"
	"time"
)

type Store interface {
	Add(tweet stream.Tweet)
	Latest() []stream.Tweet
}

type store struct {
	bucket *buckets
	mutex  *sync.RWMutex
}

func New(n int, interval time.Duration) Store {
	buckets := &buckets{
		list: make([]*bucket, n),
		curr: 0,
	}

	for i := 0; i < n; i++ {
		buckets.list[i] = &bucket{}
	}

	s := &store{
		bucket: buckets,
		mutex:  &sync.RWMutex{},
	}

	go s.swap(interval)
	return s
}

func (s store) swap(interval time.Duration) {
	for {
		<-time.After(interval)
		s.mutex.Lock()
		s.bucket.Prev()
		s.bucket.Clear()
		s.mutex.Unlock()
	}
}

func (s store) Add(tweet stream.Tweet) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.bucket.Add(tweet)
}

func (s store) Latest() []stream.Tweet {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.bucket.List()
}
