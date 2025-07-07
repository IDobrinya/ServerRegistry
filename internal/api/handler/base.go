package handler

import (
	"encoding/json"
	"net/http"
	"server-registry/internal/api/dto"
	"server-registry/internal/database/repositories"

	"github.com/jmoiron/sqlx"
)

type Handler struct {
	serverRepo *repositories.ServerRepository
	userRepo   *repositories.UserRepository
}

func NewHandler(db *sqlx.DB) *Handler {
	return &Handler{
		serverRepo: repositories.NewServerRepository(db),
		userRepo:   repositories.NewUserRepository(db),
	}
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		return
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, dto.ErrorResponse{Error: message})
}
