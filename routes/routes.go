package routes

import (
	"attendance-app/models"  // Make sure you're importing your models package correctly
	"github.com/gorilla/mux"
)

// InitRoutes initializes and returns the router with defined routes
func InitRoutes() *mux.Router {
	r := mux.NewRouter()

	// Example route for attendance
	r.HandleFunc("/api/attendance", models.AttendanceHandler).Methods("POST")
	r.HandleFunc("/api/", models.AttendanceHandler).Methods("GET")

	// You can add more routes here as needed

	return r
}
