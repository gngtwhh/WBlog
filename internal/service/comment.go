package service

import (
	"log/slog"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/repository"
	"github.com/gngtwhh/WBlog/pkg/sensitive"
)

type CommentService struct {
	repo     repository.CommentRepository
	acFilter *sensitive.ACFilter
	log      *slog.Logger
}

func NewCommentService(repo repository.CommentRepository, acFilter *sensitive.ACFilter, logger *slog.Logger) *CommentService {
	return &CommentService{
		repo:     repo,
		acFilter: acFilter,
		log:      logger.With("component", "comment_service"),
	}
}

func (s *CommentService) Create(comment *model.Comment) error {
	// TODO: should send err to frontend
	comment.Content = s.acFilter.Filter(comment.Content)
	if err := s.repo.Create(comment); err != nil {
		s.log.Error("failed to create comment",
			"uid", comment.UserID, "articleid", comment.ArticleID, "err", err)
		return err
	}
	return nil
}

func (s *CommentService) List(articleID int64, limit, offset int) ([]*model.Comment, error) {
	comments, err := s.repo.ListByArticleID(articleID, limit, offset)
	if err != nil {
		s.log.Error("failed to list articles", "err", err)
		return nil, err
	}
	return comments, nil
}
