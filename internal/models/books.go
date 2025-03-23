package models

import "time"

type Book struct {
	ID        int64
	Title     string
	CreatedAt time.Time
	Published int
	Pages     int
	Genres    []string
	Version   int
}
