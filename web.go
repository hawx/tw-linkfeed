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
	"os"
	"time"
)

var (
	consumerKey    = flag.String("consumer-key", os.Getenv("TWITTER_CONSUMER_KEY"), "")
	consumerSecret = flag.String("consumer-secret", os.Getenv("TWITTER_CONSUMER_SECRET"), "")
	accessToken    = flag.String("access-token", os.Getenv("TWITTER_OAUTH_TOKEN"), "")
	accessSecret   = flag.String("access-secret", os.Getenv("TWITTER_OAUTH_TOKEN_SECRET"), "")
)

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

		fmt.Fprintf(w, rss)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
