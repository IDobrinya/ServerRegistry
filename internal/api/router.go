package api

import (
	"server-registry/internal/api/handler"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *mux.Router {
	h := handler.NewHandler(db)

	router := mux.NewRouter()

	api := router.PathPrefix("/api/v1").Subrouter()

	// Server endpoints
	api.HandleFunc("/servers/register", h.RegisterServer).Methods("POST")
	api.HandleFunc("/servers/tunnel", h.SetServerTunnel).Methods("PUT")
	api.HandleFunc("/servers/devices", h.UnlinkAllDevices).Methods("DELETE")

	// User endpoints
	api.HandleFunc("/user/server", h.GetUserServer).Methods("GET")
	api.HandleFunc("/user/link-server", h.LinkServerToUser).Methods("POST")

	return router
}
