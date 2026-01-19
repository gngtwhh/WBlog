package service

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
)

var (
	ErrArticleNotFound = errors.New("article not found")
)

type ArticleService struct {
	repo repository.ArticleRepository
	log  *slog.Logger
}

func NewArticleService(repo repository.ArticleRepository, logger *slog.Logger) *ArticleService {
	return &ArticleService{
		repo: repo,
		log:  logger.With("componend", "article_service"),
	}
}

func (svc *ArticleService) ListArticles(limit, offset int) ([]model.Article, error) {
	articles, err := svc.repo.GetList(limit, offset)
	if err != nil {
		svc.log.Error("failed to list articles", "err", err)
		return nil, err
	}
	return articles, nil
}

func (svc *ArticleService) Count() (int64, error) {
	count, err := svc.repo.Count()
	if err != nil {
		svc.log.Error("failed to count articles", "err", err)
		return 0, err
	}
	return count, nil
}

func (svc *ArticleService) GetArticle(id int64) (model.Article, error) {
	article, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Article{}, ErrArticleNotFound
		}
		svc.log.Error("failed to get article", "id", id, "err", err)
		return model.Article{}, err
	}
	return article, nil
}

func (svc *ArticleService) Create(article *model.Article) error {
	if article.Abstract == "" {
		article.Abstract = article.Content[:100]
	}
	err := svc.repo.Create(article)
	if err != nil {
		svc.log.Error("failed to create article", "title", article.Title, "err", err)
		return err
	}
	return nil
}

func (svc *ArticleService) Update(article *model.Article) error {
	err := svc.repo.Update(article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrArticleNotFound
		}
		svc.log.Error("failed to update article", "id", article.ID, "err", err)
		return err
	}
	return nil
}

func (svc *ArticleService) Delete(id int64) error {
	err := svc.repo.Delete(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrArticleNotFound
		}
		svc.log.Error("failed to delete article", "id", id, "err", err)
		return err
	}
	return nil
}
