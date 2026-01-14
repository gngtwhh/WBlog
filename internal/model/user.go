package model

import "time"

const (
	RoleUser  = 1  // comment, view article.
	RoleAdmin = 99 // publish, edit, delete article and manage users.

	StatusNormal = 1 // functioning normally.
	StatusBanned = 0 // have been banned.
)

type User struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"` // must unique
	Password string `json:"-"`        // bcrypt hashed

	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`

	Role      int       `json:"role"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}
