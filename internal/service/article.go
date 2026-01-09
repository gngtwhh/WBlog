package service

import (
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
)

type ArticleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) *ArticleService {
	return &ArticleService{repo: repo}
}

func (h *ArticleService) ListArticles(limit, offset int) ([]model.Article, error) {
	return h.repo.GetList(limit, offset)
}

func (h *ArticleService) Count() (int64, error) {
	return h.repo.Count()
}

func (h *ArticleService) GetArticle(id int64) (model.Article, error) {
	return h.repo.GetByID(id)
}

func (h *ArticleService) Create(article *model.Article) error {
	return h.repo.Create(article)
}

func (h *ArticleService) Update(article *model.Article) error {
	return h.repo.Update(article)
}

func (h *ArticleService) Delete(id int64) error {
	return h.repo.Delete(id)
}
