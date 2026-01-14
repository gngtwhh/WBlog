package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/gngtwhh/WBlog/internal/model"
)

// ArticleRepo implements the repository.ArticleRepository interface.
type ArticleRepo struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// Create inserts a new article into the database.
// article.ID will be set if Create success.
func (r *ArticleRepo) Create(article *model.Article) error {
	query := `
		INSERT INTO articles (title,author,content,abstract,view_count)
		VALUES (?,?,?,?,?)
	`
	result, err := r.db.Exec(query, article.Title, article.Author, article.Content, article.Abstract,
		article.ViewCount)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	article.ID = uint64(id)
	return nil
}

func (r *ArticleRepo) GetByID(id int64) (model.Article, error) {
	query := `
			SELECT id, title, author, content, abstract, view_count, created_at, updated_at
			FROM articles
			WHERE id = ?
		`
	var a model.Article
	err := r.db.QueryRow(query, id).Scan(
		&a.ID, &a.Title, &a.Author, &a.Content, &a.Abstract,
		&a.ViewCount, &a.CreatedAt, &a.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Article{}, errors.New("article not found")
		}
		return model.Article{}, err
	}
	return a, nil
}

func (r *ArticleRepo) Update(article *model.Article) error {
	query := `
		UPDATE articles
		SET title=?, author=?, content=?, abstract=?
		WHERE id=?
	`
	res, err := r.db.Exec(query,
		article.Title,
		article.Author,
		article.Content,
		article.Abstract,
		article.ID,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("article with id %d not found", article.ID)
	}
	return nil
}

func (r *ArticleRepo) Delete(id int64) error {
	query := "DELETE FROM articles WHERE id = ?"
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("article with id %d not found", id)
	}
	return nil
}

// GetList retrieves a list of articles from the database.
func (r *ArticleRepo) GetList(limit, offset int) ([]model.Article, error) {
	if limit <= 0 {
		limit = 10
	}
	if offset <= 0 {
		offset = 0
	}

	query := `
		SELECT id, title, author, abstract, view_count, created_at, updated_at
		FROM articles
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]model.Article, 0, limit)
	for rows.Next() {
		var article model.Article
		if err := rows.Scan(&article.ID, &article.Title, &article.Author, &article.Abstract,
			&article.ViewCount, &article.CreatedAt, &article.UpdatedAt); err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return articles, nil
}
func (r *ArticleRepo) Count() (int64, error) {
	var count int64
	query := "SELECT count(*) FROM articles"
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
