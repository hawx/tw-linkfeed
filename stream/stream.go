package stream

import (
	"github.com/garyburd/go-oauth/oauth"
	"log"
)

type auth struct {
	Oauth       *oauth.Client
	Credentials *oauth.Credentials
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
	err := conn.Open()
	if err != nil {
		log.Fatal(err)
	}

	return conn.out
}
