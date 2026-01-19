package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/service"
	"github.com/gngtwhh/WBlog/pkg/errcode"
	"github.com/gngtwhh/WBlog/pkg/response"
)

type ArticleHandler struct {
	svc *service.ArticleService
}

// CreateArticleRequest bind POST request data, and will be cleaned to match model.Article
type CreateArticleRequest struct {
	Title    string `json:"title"`
	Author   string `json:"author"`
	Content  string `json:"content"`
	Abstract string `json:"abstract"`
}

// UpdateArticleRequest bind POST request data, and will be cleaned to match model.Article
type UpdateArticleRequest struct {
	ID uint64 `json:"id"`
	CreateArticleRequest
}

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

// ListArticles handle a GET request to gen a list of Atricles
// GET req requires two params:
// @page: page index(start from 1)
// @pagesize: count of articles per-page
func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	pageSize := r.URL.Query().Get("pagesize")
	page := r.URL.Query().Get("page")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		response.Fail(w, errcode.ParamError, "Invalid param: pagesize")
		// http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		response.Fail(w, errcode.ParamError, "Invalid param: page")
		// http.Error(w, "Invalid page", http.StatusBadRequest)
		return
	}

	if pageSizeInt <= 0 {
		pageSizeInt = 10 // default 10 articles per-page
	}
	if pageInt <= 0 {
		pageInt = 1
	}

	offset := (pageInt - 1) * pageSizeInt
	articles, err := h.svc.ListArticles(pageSizeInt, offset)
	if err != nil {
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, articles)
}

// Count handles GET req, and returns the total number of articles.
func (h *ArticleHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.Count()
	if err != nil {
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, map[string]int64{"count": count})
}

// GetArticle returns an article by id.
// GET req requires one param:
// @id: id of article required
func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(w, errcode.ParamError, "Invalid param: id")
		// http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.svc.GetArticle(int64(id))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Fail(w, errcode.ArticleNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, article)
}

// Create handles POST reqs, and create a new article.
// POST data must bind to CreateArticleRequest.
func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		// http.Error(w, "Invalid json request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		response.Fail(w, errcode.ParamError, "Title and content must not be empty")
		// http.Error(w, "Title and content must not be empty", http.StatusBadRequest)
		return
	}

	article := model.Article{
		Title:    req.Title,
		Author:   req.Author,
		Content:  req.Content,
		Abstract: req.Abstract,
	}
	err := h.svc.Create(&article)
	if err != nil {
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, map[string]uint64{"id": article.ID})
}

// Update handles POST reqs, and update the article.
// POST data must bind to UpdateArticleRequest.
func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req UpdateArticleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		// http.Error(w, "Invalid json request body", http.StatusBadRequest)
		return
	}

	if req.ID == 0 {
		response.Fail(w, errcode.ParamError, "Need ID of the article")
		// http.Error(w, "Need ID of the article", http.StatusBadRequest)
	}
	// // update allow empty content?
	// if req.Title == "" || req.Content == "" {
	// 	http.Error(w, "Title and content must not be empty", http.StatusBadRequest)
	// 	return
	// }

	article := model.Article{
		ID:       req.ID,
		Title:    req.Title,
		Author:   req.Author,
		Content:  req.Content,
		Abstract: req.Abstract,
	}
	err := h.svc.Update(&article)
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Fail(w, errcode.ArticleNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, map[string]uint64{"id": article.ID})
}

// DELETE req requires one param:
// @id: id of article required
func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(w, errcode.ParamError, "Invalid param: id")
		// http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	err = h.svc.Delete(int64(idInt))
	if err != nil {
		if errors.Is(err, service.ErrArticleNotFound) {
			response.Fail(w, errcode.ArticleNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		// http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	response.Success(w, nil)
}
