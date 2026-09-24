package http

import (
	"context"
	"net/http"
	"strconv"
)

type contextKey string

const UserIdCtxKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		UserIDStr := r.Header.Get("X-User-ID")
		if UserIDStr == "" {
			http.Error(w, "Missing X-User-Id header", http.StatusUnauthorized)
			return
		}
		UserID, err := strconv.ParseInt(UserIDStr, 10, 64)
		if err != nil || UserID <= 0 {
			http.Error(w, "Invalid UserId format", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIdCtxKey, UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
