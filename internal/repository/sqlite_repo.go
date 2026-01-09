package repository

import (
	"database/sql"
	"errors"

	"github.com/gngtwhh/WBlog/internal/model"
)

// ArticleRepo implements the repository.ArticleRepository interface.
type ArticleRepo struct {
	db *sql.DB
}

func NewArticleRepo(db *sql.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

// GetList retrieves a list of articles from the database.
func (r *ArticleRepo) GetList(count int) ([]model.Article, error) {
	if count <= 0 {
		return nil, errors.New("count must gt 0")
	}

	query := `
		SELECT id, title, author, abstract, view_count, created_at, updated_at
		FROM articles
		ORDER BY created_at DESC
		LIMIT ?
	`
	rows, err := r.db.Query(query, count)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := make([]model.Article, 0, count)
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

// Create inserts a new article into the database.
// article.ID will be set if Create success.
func (r *ArticleRepo) Create(article *model.Article) error {
	query := `
		INSERT INTO articles (title,author,content,abstract,view_count,created_at,updated_at)
		VALUES (?,?,?,?,?,?,?)
	`
	result, err := r.db.Exec(query, article.Title, article.Author, article.Content, article.Abstract,
		article.ViewCount, article.CreatedAt, article.UpdatedAt)
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
