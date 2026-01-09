package repository

import "github.com/gngtwhh/WBlog/internal/model"

// ArticleRepository defines the methods for interacting with articles in the repository.
type ArticleRepository interface {
	// Single article
	Create(article *model.Article) error
	GetByID(id int64) (model.Article, error)
	Update(article *model.Article) error
	Delete(id int64) error
	// list
	GetList(limit, offset int) ([]model.Article, error)
	Count() (int64, error)
}
