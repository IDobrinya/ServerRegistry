package handler

import (
	"encoding/json"
	"net/http"
	"server-registry/internal/api/dto"
	"strings"
)

func (h *Handler) GetUserServer(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil || !user.LinkedServer.Valid {
		writeJSON(w, http.StatusOK, dto.GetUserServerResponse{
			BridgeURL: nil,
		})
		return
	}

	server, err := h.serverRepo.GetServerByID(int(user.LinkedServer.Int32))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get server")
		return
	}

	if server == nil {
		writeJSON(w, http.StatusOK, dto.GetUserServerResponse{
			BridgeURL: nil,
		})
		return
	}

	var bridgeURL *string
	if server.BridgeURL.Valid {
		bridgeURL = &server.BridgeURL.String
	}

	writeJSON(w, http.StatusOK, dto.GetUserServerResponse{
		BridgeURL: bridgeURL,
	})
}

func (h *Handler) LinkServerToUser(w http.ResponseWriter, r *http.Request) {
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

	var req dto.LinkServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.ServerToken == "" {
		writeError(w, http.StatusBadRequest, "Server token is required")
		return
	}

	server, err := h.serverRepo.GetServerByToken(req.ServerToken)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get server")
		return
	}

	if server == nil {
		writeError(w, http.StatusNotFound, "Server not found")
		return
	}

	user, err := h.userRepo.GetUserByID(userID)
	if user == nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	_, err = h.userRepo.UpdateUserLinkedServer(userID, server.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Server linked successfully"})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
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

	existingUser, err := h.userRepo.GetUserByID(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check existing user")
		return
	}

	if existingUser != nil {
		writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "User already exists"})
		return
	}

	_, err = h.userRepo.CreateUser(userID, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, dto.SuccessResponse{Message: "User created successfully"})
}

func (h *Handler) UnlinkServer(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.userRepo.GetUserByID(userID)
	if user == nil {
		writeError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	_, err = h.userRepo.UnlinkUserDevice(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}
	writeJSON(w, http.StatusCreated, dto.SuccessResponse{Message: "Server unlinked successfully"})
}
