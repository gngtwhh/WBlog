package repository

import (
	"database/sql"
	"fmt"

	"github.com/gngtwhh/WBlog/internal/model"
)

// UserRepo implements the repository.UserRepository interface.
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(user *model.User) error {
	query := `
		INSERT INTO users (username, password,nickname,avatar,role,status)
		VALUES (?,?,?,?,?,?)
	`
	result, err := r.db.Exec(query,
		user.Username,
		user.Password,
		user.Nickname,
		user.Avatar,
		user.Role,
		user.Status,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}

func (r *UserRepo) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, password, nickname, avatar, role, status, created_at, updated_at
		FROM users
		WHERE username = ?
	`
	row := r.db.QueryRow(query, username)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) GetByID(id uint64) (*model.User, error) {
	query := `
		SELECT id, username, password, nickname, avatar, role, status, created_at, updated_at
		FROM users
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)
	user := &model.User{}
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Nickname,
		&user.Avatar,
		&user.Role,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepo) Update(user *model.User) error {
	query := `
			UPDATE users
			SET password=?, nickname=?, avatar=?, role=?, status=?, updated_at=CURRENT_TIMESTAMP
			WHERE id=?
		`

	res, err := r.db.Exec(query,
		user.Password,
		user.Nickname,
		user.Avatar,
		user.Role,
		user.Status,
		user.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("user with id %d not found", user.ID)
	}
	return nil
}
