package main

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
	"github.com/hawx/tw-linkfeed/store"
	"github.com/hawx/tw-stream"
	"github.com/hawx/tw-linkfeed/views"
	"github.com/hawx/serve"

	"flag"
	"fmt"
	"net/http"
	"time"
)

var (
	consumerKey    = flag.String("consumer-key", "", "")
	consumerSecret = flag.String("consumer-secret", "", "")
	accessToken    = flag.String("access-token", "", "")
	accessSecret   = flag.String("access-secret", "", "")

	title = flag.String("title", "tw-linkfeed", "")
	url   = flag.String("url", "http://localhost:8080/", "")

	port   = flag.String("port", "8080", "")
	socket = flag.String("socket", "", "")
	help   = flag.Bool("help", false, "")
)

const HELP = `Usage: tw-linkfeed [options]

  Serves a feed (in html at '/', and rss at '/feed') of all links in
  your twitter timeline.

    --consumer-key <value>
    --consumer-secret <value>
    --access-token <value>
    --access-secret <value>

    --title <title>     # Title of page/feed (default: 'tw-linkfeed')
    --url <url>         # URL running at (default: 'http://localhost:8080/')

    --port <port>       # Port to run on (default: '8080')
    --socket <path>     # Serve using a unix socket instead
    --help              # Display this help message
`

func run(store store.Store) {
	auth := stream.Auth(*consumerKey, *consumerSecret, *accessToken, *accessSecret)

	for tweet := range stream.Timeline(auth) {
		if len(tweet.Entities.Urls) > 0 {
			go func(tweet stream.Tweet) {
				doc, err := goquery.NewDocument(*tweet.Entities.Urls[0].ExpandedUrl)
				if err != nil {
					store.Add(tweet)
					return
				}

				title := doc.Find("title").Text()
				if title != "" {
					tweet.Entities.Urls[0].DisplayUrl = title
				}
				store.Add(tweet)
			}(tweet)
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
		views.List.Execute(w, struct {
			Tweets []stream.Tweet
			Url    string
			Title  string
		}{store.Latest(), *url, *title})
	})

	http.HandleFunc("/feed", func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		feed := &feeds.Feed{
			Title:   *title,
			Link:    &feeds.Link{Href: *url + "feed"},
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

		w.Header().Add("Content-Type", "application/rss+xml")

		err := feed.WriteRss(w)
		if err != nil {
			w.WriteHeader(500)
			return
		}
	})

	serve.Serve(*port, *socket, http.DefaultServeMux)
}
