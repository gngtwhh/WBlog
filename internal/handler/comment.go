package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gngtwhh/WBlog/internal/middleware"
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/service"
	"github.com/gngtwhh/WBlog/pkg/errcode"
	"github.com/gngtwhh/WBlog/pkg/response"
)

type CommentHandler struct {
	commentsvc *service.CommentService
	articlesvc *service.ArticleService
}

type CreateCommentReq struct {
	ArticleID int64  `json:"article_id"`
	Content   string `json:"content"`
}

func NewCommentHandler(commentsvc *service.CommentService, articlesvc *service.ArticleService) *CommentHandler {
	return &CommentHandler{
		commentsvc: commentsvc,
		articlesvc: articlesvc,
	}
}

func (h *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
		return
	}

	var req CreateCommentReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		return
	}

	if req.ArticleID <= 0 {
		response.Fail(w, errcode.ParamError, "Invalid article ID")
		return
	}
	if req.Content == "" {
		response.Fail(w, errcode.ParamError, "Comment content cannot be empty")
		return
	}

	if _, err := h.articlesvc.GetArticle(req.ArticleID); err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Fail(w, errcode.ArticleNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
	}

	username, _ := middleware.GetUsername(r)
	comment := &model.Comment{
		UserID:    userID,
		ArticleID: uint64(req.ArticleID),
		Username:  username,
		Content:   req.Content,
	}

	if err := h.commentsvc.Create(comment); err != nil {
		response.Fail(w, errcode.ServerError)
		return
	}
	response.Success(w, map[string]uint64{"id": comment.ID})
}

func (h *CommentHandler) ListComments(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	articleIDStr := query.Get("article_id")
	articleID, err := strconv.ParseInt(articleIDStr, 10, 64)
	if err != nil || articleID <= 0 {
		response.Fail(w, errcode.ParamError, "Invalid or missing article_id")
		return
	}

	pageStr := query.Get("page")
	pageSizeStr := query.Get("page_size")
	page, _ := strconv.Atoi(pageStr)
	if page <= 0 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(pageSizeStr)
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	} else if pageSize > 100 {
		pageSize = 100 // Max limit
	}

	offset := (page - 1) * pageSize
	comments, err := h.commentsvc.List(articleID, pageSize, offset)
	if err != nil {
		response.Fail(w, errcode.ServerError)
		return
	}

	if comments == nil {
		comments = []*model.Comment{}
	}
	response.Success(w, comments)
}
