package model

import "time"

// Article represents a blog article.
type Article struct {
	ID uint64 `json:"id"`

	Title    string `json:"title"`
	Author   string `json:"author"`
	Content  string `json:"content"`
	Abstract string `json:"abstract"`

	ViewCount uint64    `json:"view_count"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
