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

func (svc *ArticleService) ListArticles(limit, offset int) ([]model.Article, error) {
	return svc.repo.GetList(limit, offset)
}

func (svc *ArticleService) Count() (int64, error) {
	return svc.repo.Count()
}

func (svc *ArticleService) GetArticle(id int64) (model.Article, error) {
	return svc.repo.GetByID(id)
}

func (svc *ArticleService) Create(article *model.Article) error {
	return svc.repo.Create(article)
}

func (svc *ArticleService) Update(article *model.Article) error {
	return svc.repo.Update(article)
}

func (svc *ArticleService) Delete(id int64) error {
	return svc.repo.Delete(id)
}
