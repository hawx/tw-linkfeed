package stream

import (
	"github.com/garyburd/go-oauth/oauth"
	"log"
)

const (
	SAMPLE_URL   = "https://stream.twitter.com/1.1/statuses/sample.json"
	STREAM_URL   = "https://userstream.twitter.com/1.1/user.json"
	USER_URL     = STREAM_URL + "?with=user"
)

type auth struct {
	Oauth       *oauth.Client
	Credentials *oauth.Credentials
}

func (a *auth) Name() string {
	// get name of user!?!?
	return "me"
}

func Auth(consumerKey, consumerSecret, accessToken, accessSecret string) *auth {
	return &auth{
		Oauth: &oauth.Client{
			Credentials: oauth.Credentials{
				Token:  consumerKey,
				Secret: consumerSecret,
			},
		},
		Credentials: &oauth.Credentials{
			Token:  accessToken,
			Secret: accessSecret,
		},
	}
}

type Stream chan Tweet

func Timeline(creds *auth) Stream {
	conn := newConnection(creds)
	err := conn.Open(STREAM_URL)
	if err != nil {
		log.Fatal(err)
	}

	return conn.out
}

func Self(creds *auth) Stream {
	conn := newConnection(creds)
	err := conn.Open(USER_URL)
	if err != nil {
		log.Fatal(err)
	}

	return conn.out
}
