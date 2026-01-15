package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gngtwhh/WBlog/pkg/utils"
)

type ContextKey string

const UserIDKey ContextKey = "user_id"

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized request: missing token", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "unauthorized request: the token format is incorrect", http.StatusUnauthorized)
			return
		}
		tokenStr := parts[1]

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "unauthorized request: invalid or expired token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		next(w, r.WithContext(ctx))
	}
}

func GetUserID(r *http.Request) (uint64, bool) {
	id, ok := r.Context().Value(UserIDKey).(uint64)
	return id, ok
}
