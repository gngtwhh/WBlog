package service

import (
	"log/slog"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
)

type ArticleService struct {
	repo repository.ArticleRepository
	log  *slog.Logger
}

func NewArticleService(repo repository.ArticleRepository, logger *slog.Logger) *ArticleService {
	return &ArticleService{
		repo: repo,
		log:  logger,
	}
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
