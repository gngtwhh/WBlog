package repository

import (
	"database/sql"
	"log/slog"

	"github.com/gngtwhh/WBlog/internal/model"
)

// CommentRepo implements the repository.CommentRepository interface.
type CommentRepo struct {
	db  *sql.DB
	log *slog.Logger
}

func NewCommentRepo(db *sql.DB, log *slog.Logger) *CommentRepo {
	return &CommentRepo{
		db:  db,
		log: log.With("component", "comment_repo"),
	}
}

func (r *CommentRepo) Create(comment *model.Comment) error {
	query := `
		INSERT INTO comments (user_id, article_id, content, username)
		VALUES (?, ?, ?, ?)
	`
	res, err := r.db.Exec(query, comment.UserID, comment.ArticleID, comment.Content, comment.Username)
	if err != nil {
		r.log.Error("Create comment failed", slog.String("err", err.Error()))
		return err
	}

	id, _ := res.LastInsertId()
	comment.ID = uint64(id)
	return nil
}

func (r *CommentRepo) ListByArticleID(articleID int64, limit, offset int) ([]*model.Comment, error) {
	query := `
		SELECT id, user_id, article_id, content, username, created_at
		FROM comments
		WHERE article_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.Query(query, articleID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]*model.Comment, 0)
	for rows.Next() {
		var c model.Comment
		err := rows.Scan(&c.ID, &c.UserID, &c.ArticleID, &c.Content, &c.Username, &c.CreatedAt)
		if err != nil {
			continue
		}
		list = append(list, &c)
	}
	return list, nil
}
