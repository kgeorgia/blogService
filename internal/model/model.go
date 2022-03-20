package model

import "time"

type Record struct {
	Id        int       `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"created_at"`
}
