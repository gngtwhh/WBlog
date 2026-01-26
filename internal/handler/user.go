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
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		// http.Error(w, "Invalid json request body", http.StatusBadRequest)
		return
	}
	if req.Username == "" || req.Password == "" {
		response.Fail(w, errcode.ParamError, "need username and password field")
		// http.Error(w, "need username and password field", http.StatusBadRequest)
		return
	}
	if req.Password != req.ConfirmPassword {
		response.Fail(w, errcode.ParamError, "the passwords entered twice must be consistent.")
		// http.Error(w, "the passwords entered twice must be consistent.", http.StatusBadRequest)
		return
	}

	user := model.User{
		Username: req.Username,
		Password: req.Password,
		Nickname: req.Nickname,
	}
	if err := h.svc.Register(&user); err != nil {
		if errors.Is(err, service.ErrUserExists) {
			response.Fail(w, errcode.UserExists)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}
	response.Success(w, map[string]interface{}{
		"id": user.ID,
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
		if errors.Is(err, service.ErrAuthFailed) {
			response.Fail(w, errcode.AuthFailed)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}
	resp := map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":       user.ID,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"role":     user.Role,
		},
	}
	response.Success(w, resp)
}

func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
	}
	user, err := h.svc.GetProfile(userID)
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(w, errcode.UserNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}
	response.Success(w, user)
}

func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		return
	}

	user := &model.User{
		ID:       userID,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}

	if err := h.svc.UpdateProfile(user); err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(w, errcode.UserNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}
	response.Success(w, nil)
}

func (h *UserHandler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
		return
	}
	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Fail(w, errcode.ParamError, "Invalid json request body")
		return
	}
	if err := h.svc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		if errors.Is(err, service.ErrInvalidOldPass) {
			response.Fail(w, errcode.AuthFailed, "Old password incorrect")
			return
		}
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(w, errcode.UserNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}
	response.Success(w, nil)
}

func (h *UserHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// TODO: need implement
	response.Fail(w, errcode.ServerError, "feature not implemented yet")
}

func (h *UserHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		response.Fail(w, errcode.ParamError, "invalid user id")
		return
	}

	user, err := h.svc.GetProfile(uint64(id))
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			response.Fail(w, errcode.UserNotFound)
			return
		}
		response.Fail(w, errcode.ServerError)
		return
	}

	publicInfo := map[string]interface{}{
		"id":         user.ID,
		"nickname":   user.Nickname,
		"avatar":     user.Avatar,
		"created_at": user.CreatedAt,
	}

	response.Success(w, publicInfo)
}

func (h *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	tokenStr, ok := middleware.GetTokenRaw(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
		return
	}

	exp, ok := middleware.GetClaimsExp(r)
	if !ok {
		response.Fail(w, errcode.TokenInvalid)
		return
	}

	if err := h.svc.Logout(tokenStr, exp); err != nil {
		response.Fail(w, errcode.ServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{Name: "token", MaxAge: -1, Path: "/"})
	response.Success(w, nil)
}
