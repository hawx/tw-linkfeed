package stream

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	DIAL_TIMEOUT = 5 * time.Second
)

type conn struct {
	client  *http.Client
	auth    *auth
	out     chan Tweet
}

func newConnection(creds *auth) *conn {
	client := &http.Client{}
	out := make(chan Tweet)

	return &conn{client: client, out: out, auth: creds}
}

func (c conn) Open(streamUrl string) error {
	req, _ := http.NewRequest("GET", streamUrl, nil)
	req.Header.Set("Authorization", c.auth.Oauth.AuthorizationHeader(c.auth.Credentials, "GET", req.URL, nil))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("stream: %s", err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("stream: (%d) %s", resp.StatusCode, body)
	}

	go c.run(resp.Body)

	return nil
}

func (c conn) run(body io.Reader) {
	decoder := json.NewDecoder(body)

	for {
		var tweet Tweet
		if err := decoder.Decode(&tweet); err != nil {
			log.Fatal(err)
		}

		if tweet.Id != 0 {
			c.out <- tweet
		}
	}
}
