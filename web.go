package main

import (
	"github.com/gorilla/feeds"
	"github.com/hawx/tw-linkfeed/store"
	"github.com/hawx/tw-linkfeed/stream"
	"github.com/hawx/tw-linkfeed/views"

	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	port           = flag.String("port", "8080", "")
	consumerKey    = flag.String("consumer-key", "", "")
	consumerSecret = flag.String("consumer-secret", "", "")
	accessToken    = flag.String("access-token", "", "")
	accessSecret   = flag.String("access-secret", "", "")
	help           = flag.Bool("help", false, "")
)

const HELP = `Usage: tw-linkfeed [options]

  Serves a feed (in html at '/', and rss at '/feed') of all links in
  your twitter timeline.

    --port <port>       # Port to run on (default: '8080')

    --consumer-key <value>
    --consumer-secret <value>
    --access-token <value>
    --access-secret <value>

    --help              # Display this help message
`

func run(store store.Store) {
	auth := stream.Auth(*consumerKey, *consumerSecret, *accessToken, *accessSecret)

	for tweet := range stream.Timeline(auth) {
		if len(tweet.Entities.Urls) > 0 {
			store.Add(tweet)
		}
	}
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(HELP)
		return
	}

	store := store.New(24, time.Hour)
	go run(store)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		views.List.Execute(w, store.Latest())
	})

	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:   "tw-linkfeed",
			Link:    &feeds.Link{Href: "/feed"},
			Created: now,
		}

		for _, tweet := range store.Latest() {
			url := tweet.Entities.Urls[0]

			feed.Items = append(feed.Items, &feeds.Item{
				Title:       url.DisplayUrl,
				Link:        &feeds.Link{Href: *url.ExpandedUrl},
				Description: tweet.Text,
				Created:     tweet.CreatedAt.Time,
			})
		}

		rss, err := feed.ToRss()
		if err != nil {
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "application/rss+xml")
		fmt.Fprintf(w, rss)
	})

	log.Println("listening on :"+*port)
	log.Fatal(http.ListenAndServe(":"+*port, nil))
}
