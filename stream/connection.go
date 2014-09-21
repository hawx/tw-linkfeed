package stream

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	SAMPLE_URL   = "https://stream.twitter.com/1.1/statuses/sample.json"
	STREAM_URL   = "https://userstream.twitter.com/1.1/user.json"
	DIAL_TIMEOUT = 5 * time.Second
)

type conn struct {
	client  *http.Client
	auth    *auth
	out     chan Tweet
	closer  io.Closer
	decoder *json.Decoder
	timeout time.Duration
}

func newConnection(creds *auth, timeout time.Duration) *conn {
	client := &http.Client{}
	out := make(chan Tweet)

	return &conn{client: client, out: out, auth: creds, timeout: timeout}
}

func (c conn) Open() error {
	req, _ := http.NewRequest("GET", STREAM_URL, nil)
	req.Header.Set("Authorization", c.auth.Oauth.AuthorizationHeader(c.auth.Credentials, "GET", req.URL, nil))

	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("stream: making filter request failed: %s", err)
	}

	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		return fmt.Errorf("stream: filter failed (%d): %s", resp.StatusCode, body)
	}

	c.closer = resp.Body
	c.decoder = json.NewDecoder(resp.Body)

	// go plain(resp.Body)

	go c.run(resp.Body)

	return nil
}

func plain(body io.Reader) {
	for {
		reader := bufio.NewReader(body)
		line, _ := reader.ReadBytes('\r')
		log.Println(string(line))
	}
}

func (c conn) run(body io.Reader) {
	decoder := json.NewDecoder(body)

	for {
		var tweet Tweet
		if c.timeout != 0 {
			c.client.Timeout = c.timeout
		}

		if err := decoder.Decode(&tweet); err != nil {
			log.Fatal(err)
		}

		if tweet.Id != 0 {
			c.out <- tweet
		}
	}
}

func (c conn) Close() error {
	// if err := c.conn.Close(); err != nil {
	// 	c.closer.Close()
	// 	return err
	// }
	return c.closer.Close()
}
