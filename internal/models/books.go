package models

import "time"

type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"-"`
	Published int       `json:"published"`
	Pages     int       `json:"pages,omitempty,string"`
	Genres    []string  `json:"genres,omitempty"`
	Version   int       `json:"version"`
}
