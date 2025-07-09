package handler

import (
	"encoding/json"
	"net/http"
	"server-registry/internal/api/dto"
)

func (h *Handler) RegisterServer(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterServerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "Server token is required")
		return
	}

	existingServer, err := h.serverRepo.GetServerByToken(req.Token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check existing server")
		return
	}

	if existingServer != nil {
		writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Server already registered"})
		return
	}

	_, err = h.serverRepo.CreateServer(req.Token, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to register server")
		return
	}

	writeJSON(w, http.StatusCreated, dto.SuccessResponse{Message: "Server registered successfully"})
}

func (h *Handler) SetServerTunnel(w http.ResponseWriter, r *http.Request) {
	var req dto.SetTunnelRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "Server token is required")
		return
	}

	if req.BridgeURL == "" {
		writeError(w, http.StatusBadRequest, "Bridge URL is required")
		return
	}

	existingServer, err := h.serverRepo.GetServerByToken(req.Token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check server")
		return
	}

	if existingServer == nil {
		writeError(w, http.StatusNotFound, "Server not found")
		return
	}

	_, err = h.serverRepo.UpdateServerURL(req.Token, req.BridgeURL)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update server tunnel")
		return
	}

	writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "Server tunnel updated successfully"})
}

func (h *Handler) UnlinkAllDevices(w http.ResponseWriter, r *http.Request) {
	var req dto.UnlinkDevicesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "Server token is required")
		return
	}

	existingServer, err := h.serverRepo.GetServerByToken(req.Token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to check server")
		return
	}

	if existingServer == nil {
		writeError(w, http.StatusNotFound, "Server not found")
		return
	}

	rowsAffected, err := h.userRepo.UnlinkAllUsersFromServer(req.Token)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to unlink devices")
		return
	}

	if rowsAffected == 0 {
		writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "No devices were connected to unlink"})
	} else {
		writeJSON(w, http.StatusOK, dto.SuccessResponse{Message: "All devices unlinked successfully"})
	}
}
