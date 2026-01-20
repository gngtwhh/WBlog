package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gngtwhh/WBlog/pkg/errcode"
	"github.com/gngtwhh/WBlog/pkg/response"
	"github.com/gngtwhh/WBlog/pkg/utils"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"
const UsernameKey ContextKey = "username"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: missing token")
			// http.Error(w, "unauthorized request: missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: the token format is incorrect")
			return
		}
		tokenStr := parts[1]

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		next(w, r.WithContext(ctx))
	}
}

func GetUserID(r *http.Request) (uint64, bool) {
	id, ok := r.Context().Value(UserIDKey).(uint64)
	return id, ok
}

func GetUsername(r *http.Request) (string, bool) {
	username, ok := r.Context().Value(UsernameKey).(string)
	return username, ok
}
