package stream

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	TimeFormat = "Mon Jan _2 15:04:05 +0000 2006"
)

type Time struct {
	time.Time
}

func (t Time) MarshalJSON() ([]byte, error) {
	formatted := t.Format(TimeFormat)
	return json.Marshal(formatted)
}

func (t *Time) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	t.Time, err = time.Parse(TimeFormat, s)
	return nil
}

type Url struct {
	Url         string  `json:"url"`
	DisplayUrl  string  `json:"display_url"`
	ExpandedUrl *string `json:"expanded_url"`
}

type Entities struct {
	Urls []Url `json:"urls"`
}

type User struct {
	Id         int64   `json:"id"`
	Name       string  `json:"name"`
	ScreenName string  `json:"screen_name"`
	Url        *string `json:"url"`
}

type Tweet struct {
	Id        int64    `json:"id"`
	Source    string   `json:"source"`
	CreatedAt Time     `json:"created_at"`
	Entities  Entities `json:"entities"`
	Text      string   `json:"text"`
	User      User     `json:"user"`
}

func (t Tweet) Link() string {
	return "https://twitter.com/" + t.User.ScreenName + "/status/" + strconv.FormatInt(t.Id, 10)
}
