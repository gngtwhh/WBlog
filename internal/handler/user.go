package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gngtwhh/WBlog/internal/middleware"
	"github.com/gngtwhh/WBlog/internal/model"
	"github.com/gngtwhh/WBlog/internal/service"
)

type UserHandler struct {
	svc *service.UserService
}

type RegisterRequest struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	Nickname        string `json:"nickname"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid json request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		http.Error(w, "need username and password field", http.StatusBadRequest)
		return
	}
	if req.Password != req.ConfirmPassword {
		http.Error(w, "the passwords entered twice must be consistent.", http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	if err := h.svc.Register(&user); err != nil {
		// TODO: define specified error code
		http.Error(w, "register failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":  user.ID,
		"msg": "ok",
	})
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json request body", http.StatusBadRequest)
		return
	}
	user, token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		http.Error(w, "login failed: "+err.Error(), http.StatusUnauthorized)
		return
	}
	resp := map[string]interface{}{
		"token": token,
		"user":  map[string]interface{}{},
		"role":  user.Role,
	}
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unable to obtain user information", http.StatusUnauthorized)
	}
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unable to obtain user information", http.StatusUnauthorized)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json request body", http.StatusBadRequest)
		return
	}

	user := &model.User{
		ID:       userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}

	if err := h.svc.UpdateProfile(user); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"msg": "ok"})
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unable to obtain user information", http.StatusUnauthorized)
		return
	}
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json request body", http.StatusBadRequest)
		return
	}
	if err := h.svc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		http.Error(w, "update failed: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"msg": "ok"})
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// TODO: need implement
	http.Error(w, "need implement", http.StatusInternalServerError)
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid user id", http.StatusBadRequest)
		return
	}

	user, err := h.svc.GetProfile(uint64(id))
	if err != nil {
		http.Error(w, "user does not exist", http.StatusNotFound)
		return
	}

	publicInfo := map[string]interface{}{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
	}
	json.NewEncoder(w).Encode(publicInfo)
}
