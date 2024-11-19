package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"
)

// Attendance struct defines the model for the attendance data
type Attendance struct {
	ID           int    `json:"id"`
	UserName     string `json:"user_name"`
	TicketCode   string `json:"ticket_code"`
	AttendStatus string `json:"attend_status"`
	AttendTime   string `json:"attend_time"`
}

// isValidTicketCode checks if the ticket code is non-empty and valid
func isValidTicketCode(ticketCode string) bool {
	// Ticket code should not be empty
	return ticketCode != ""
}

// isValidAttendTime checks if the attend time is in a valid format (e.g., YYYY-MM-DDTHH:MM:SS)
func isValidAttendTime(attendTime string) bool {
	_, err := time.Parse("2006-01-02T15:04:05", attendTime)
	return err == nil
}

// AttendanceHandler is the handler for the /api/attendance endpoint
func AttendanceHandler(w http.ResponseWriter, r *http.Request) {
	var attendance Attendance

	// Decode the JSON request body into the Attendance struct
	err := json.NewDecoder(r.Body).Decode(&attendance)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate the TicketCode
	if !isValidTicketCode(attendance.TicketCode) {
		http.Error(w, "Ticket Code is required", http.StatusBadRequest)
		return
	}

	// Validate the AttendTime format
	if !isValidAttendTime(attendance.AttendTime) {
		http.Error(w, "Invalid Attend Time format", http.StatusBadRequest)
		return
	}

	// Print the attendance data (or you can save to DB here)
	fmt.Printf("Received attendance: %+v\n", attendance)

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(attendance)
}
