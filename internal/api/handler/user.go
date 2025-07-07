package handler

import (
	"encoding/json"
	"net/http"
	"server-registry/internal/api/dto"
)

func (h *Handler) GetUserServer(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		writeError(w, http.StatusInternalServerError, "User ID not found in context")
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
		ServerToken: server.Token,
		BridgeURL:   bridgeURL,
	})
}

func (h *Handler) LinkServerToUser(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		writeError(w, http.StatusInternalServerError, "User ID not found in context")
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
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	if user == nil {
		_, err = h.userRepo.CreateUser(userID, &req.ServerToken)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}
	} else {
		_, err = h.userRepo.UpdateUserLinkedServer(userID, req.ServerToken)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Failed to update user")
			return
		}
	}

	writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Server linked successfully"})
}
