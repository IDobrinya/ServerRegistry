package handler

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const UserIDContextKey contextKey = "userID"

func (h *Handler) UserIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get("User-ID")
		if userID == "" {
			writeError(w, http.StatusBadRequest, "User-ID header is required")
			return
		}

		userID = strings.TrimSpace(userID)
		if userID == "" {
			writeError(w, http.StatusBadRequest, "User-ID header cannot be empty")
			return
		}

		ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserIDFromContext(r *http.Request) string {
	userID, ok := r.Context().Value(UserIDContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}
