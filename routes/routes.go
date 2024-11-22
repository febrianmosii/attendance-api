package routes

import (
	"attendance-app/internal/api"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Add routes
	router.HandleFunc("/api/v1/sync/{operatorId}", api.SyncHandler).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/sync", api.SyncPutHandler).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/operator/register", api.RegisterHandler).Methods(http.MethodPost)

	return router
}
