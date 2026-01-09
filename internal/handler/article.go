package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/service"
)

type ArticleHandler struct {
	svc *service.ArticleService
}

func NewArticleHandler(svc *service.ArticleService) *ArticleHandler {
	return &ArticleHandler{svc: svc}
}

func (h *ArticleHandler) ListArticles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	// get param:
	// @page: page index(start from 1)
	// @pagesize: count of articles per-page
	pageSize := r.URL.Query().Get("pagesize")
	page := r.URL.Query().Get("page")

	pageSizeInt, err := strconv.Atoi(pageSize)
	if err != nil {
		http.Error(w, "Invalid page size", http.StatusBadRequest)
		return
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		http.Error(w, "Invalid page", http.StatusBadRequest)
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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(articles)
}

func (h *ArticleHandler) Count(w http.ResponseWriter, r *http.Request) {
	count, err := h.svc.Count()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(count)
}

func (h *ArticleHandler) GetArticle(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	article, err := h.svc.GetArticle(int64(id))
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(article)
}

func (h *ArticleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var article model.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.svc.Create(&article)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(article.ID)
}

func (h *ArticleHandler) Update(w http.ResponseWriter, r *http.Request) {
	var article model.Article
	if err := json.NewDecoder(r.Body).Decode(&article); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.svc.Update(&article)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(article.ID)
}

func (h *ArticleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid article ID", http.StatusBadRequest)
		return
	}

	id := int64(idInt)
	err = h.svc.Delete(id)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(id)
}
