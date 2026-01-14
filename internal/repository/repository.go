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

// UserRepository defines the method for managing users of blog webpages.
type UserRepository interface {
	Create(user *model.User) error
	GetByUsername(username string) (*model.User, error)
	GetByID(id uint64) (*model.User, error)
	Update(user *model.User) error
}
