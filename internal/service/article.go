package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/gngtwhh/WBlog/internal/cache"
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/redis/go-redis/v9"
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
	cacheKey := fmt.Sprintf("article:detail:%d", id)
	ctx := context.Background()
	val, err := cache.RDB.Get(ctx, cacheKey).Result()
	if err == nil {
		var article model.Article
		if jsonErr := json.Unmarshal([]byte(val), &article); jsonErr == nil {
			return article, nil
		}
		svc.log.Warn("failed to unmarshal cached article", "id", id, "err", err)
	} else if err != redis.Nil {
		svc.log.Warn("redis error during get", "key", cacheKey, "err", err)
	}

	// access db
	article, err := svc.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Article{}, ErrArticleNotFound
		}
		svc.log.Error("failed to get article", "id", id, "err", err)
		return model.Article{}, err
	}
	// update cache
	data, marshalErr := json.Marshal(article)
	if marshalErr == nil {
		setErr := cache.RDB.Set(ctx, cacheKey, data, time.Hour).Err()
		if setErr != nil {
			svc.log.Warn("failed to set cache", "key", cacheKey, "err", setErr)
		}
	}
	return article, nil
}

func (svc *ArticleService) Create(article *model.Article) error {
	svc.ensureAbstract(article)
	err := svc.repo.Create(article)
	if err != nil {
		svc.log.Error("failed to create article", "title", article.Title, "err", err)
		return err
	}
	return nil
}

func (svc *ArticleService) Update(article *model.Article) error {
	svc.ensureAbstract(article)
	err := svc.repo.Update(article)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrArticleNotFound
		}
		svc.log.Error("failed to update article", "id", article.ID, "err", err)
		return err
	}
	// delete cache
	cacheKey := fmt.Sprintf("article:detail:%d", article.ID)
	if delErr := cache.RDB.Del(context.Background(), cacheKey).Err(); delErr != nil {
		svc.log.Warn("failed to delete cache after update", "key", cacheKey, "err", delErr)
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
	// delete cache
	cacheKey := fmt.Sprintf("article:detail:%d", id)
	if delErr := cache.RDB.Del(context.Background(), cacheKey).Err(); delErr != nil {
		svc.log.Warn("failed to delete cache after delete", "key", cacheKey, "err", delErr)
	}
	return nil
}

func (svc *ArticleService) ensureAbstract(article *model.Article) {
	if article.Abstract != "" {
		return
	}
	const summaryLen = 100
	contentRunes := []rune(article.Content)
	if len(contentRunes) > summaryLen {
		article.Abstract = string(contentRunes[:summaryLen]) + "..."
	} else {
		article.Abstract = article.Content
	}
}
