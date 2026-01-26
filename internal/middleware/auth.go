package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/gngtwhh/WBlog/internal/cache"
	"github.com/gngtwhh/WBlog/pkg/errcode"
	"github.com/gngtwhh/WBlog/pkg/response"
	"github.com/gngtwhh/WBlog/pkg/utils"
)

type ContextKey string

const (
	UserIDKey    ContextKey = "user_id"
	UsernameKey  ContextKey = "username"
	TokenRawKey  ContextKey = "token_raw"
	ClaimsExpKey ContextKey = "claims_exp"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: missing token")
			// http.Error(w, "unauthorized request: missing token", http.StatusUnauthorized)
			return
		}
		// get token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: the token format is incorrect")
			return
		}
		tokenStr := parts[1]

		// check blacklist
		key := cache.PrefixJWTBlacklist + tokenStr
		n, err := cache.RDB.Exists(context.Background(), key).Result()
		if err == nil && n > 0 {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: login has expired, please log in again")
			return
		}

		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			response.Fail(w, errcode.AuthFailed, "unauthorized request: invalid or expired token")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, UsernameKey, claims.Username)
		ctx = context.WithValue(ctx, TokenRawKey, tokenStr)
		if claims.ExpiresAt != nil {
			ctx = context.WithValue(ctx, ClaimsExpKey, claims.ExpiresAt.Unix())
		}
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

func GetTokenRaw(r *http.Request) (string, bool) {
	token, ok := r.Context().Value(TokenRawKey).(string)
	return token, ok
}

func GetClaimsExp(r *http.Request) (int64, bool) {
	exp, ok := r.Context().Value(ClaimsExpKey).(int64)
	return exp, ok
}
