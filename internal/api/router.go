package api

import (
	"log"
	"server-registry/internal/api/handler"
	"server-registry/internal/api/middleware"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func NewRouter(db *sqlx.DB) *mux.Router {
	h := handler.NewHandler(db)
	logger := log.New(log.Writer(), "[API] ", log.LstdFlags)

	router := mux.NewRouter()
	router.Use(middleware.LoggingMiddleware(logger))

	api := router.PathPrefix("/api/v1").Subrouter()

	// Server endpoints
	api.HandleFunc("/servers/register", h.RegisterServer).Methods("POST")
	api.HandleFunc("/servers/tunnel", h.SetServerTunnel).Methods("PUT")
	api.HandleFunc("/servers/devices", h.UnlinkAllDevices).Methods("DELETE")

	// User endpoints
	api.HandleFunc("/user/server", h.GetUserServer).Methods("GET")
	api.HandleFunc("/user/link-server", h.LinkServerToUser).Methods("POST")
	api.HandleFunc("/user/create", h.CreateUser).Methods("POST")
	api.HandleFunc("/user/unlink-server", h.UnlinkServer).Methods("POST")

	return router
}
