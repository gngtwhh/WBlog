package model

import "time"

type Comment struct {
	ID        uint64 `json:"id"`
	UserID    uint64 `json:"user_id"`
	ArticleID uint64 `json:"article_id"`
	Username  string `json:"username"`
	Content   string `json:"content"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
