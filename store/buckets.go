package store

import (
	"github.com/hawx/tw-stream"
)

type bucket []stream.Tweet

type buckets struct {
	list []*bucket
	curr int
}

func (b *buckets) Prev() {
	if b.curr == 0 {
		b.curr = len(b.list)
	}
	b.curr--
}

func (b *buckets) Add(tweet stream.Tweet) {
	*b.list[b.curr] = append(*b.list[b.curr], tweet)
}

func (b *buckets) Clear() {
	b.list[b.curr] = &bucket{}
}

func (b *buckets) List() []stream.Tweet {
	tweets := []stream.Tweet{}

	for i := b.curr; i < len(b.list)+b.curr; i++ {
		bucket := b.list[i%len(b.list)]
		for i := len(*bucket) - 1; i >= 0; i-- {
			tweets = append(tweets, (*bucket)[i])
		}
	}

	return tweets
}
