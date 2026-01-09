package repository

import "github.com/gngtwhh/WBlog/internal/model"

// ArticleRepository defines the methods for interacting with articles in the repository.
type ArticleRepository interface {
	GetList(count int) ([]model.Article, error)
	Create(article *model.Article) error
}
