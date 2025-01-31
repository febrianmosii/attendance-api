package api

import (
	"attendance-app/internal/helpers"
	"attendance-app/internal/models"
	"attendance-app/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// SyncHandler handles the GET /sync/{operatorId} request.
func SyncHandler(w http.ResponseWriter, r *http.Request) {
	// Extract operatorId from URL
	vars := mux.Vars(r)
	operatorId, _ := strconv.Atoi(vars["operatorId"])
	log.Println("Operator ID", operatorId)

	// Query the database for the attendance data related to this operator
	data, err := repository.GetAttendanceData(operatorId) // Use the repository package function

	if err != nil {
		var message string

		message = "Failed to get data"

		if err.Error() == "operator_not_found" {
			message = "Operator not found"
		}

		helpers.SetResponse(w, r, message, nil, http.StatusNotFound)
		return
	}

	if data == nil {
		data = []models.SyncResponse{}
	}

	// Respond with the fetched data
	helpers.SetResponse(w, r, "Request successful", data, http.StatusOK)
}
func SyncPutHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the request payload
	var payload models.SyncPayload

	// Decode the JSON payload into the SyncPayload struct
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		helpers.SetResponse(w, r, "Invalid request payload", nil, http.StatusBadRequest)
		log.Println(err)
		return
	}

	// Validate the input
	if len(payload.Data) == 0 {
		helpers.SetResponse(w, r, "Invoice codes cannot be empty", nil, http.StatusBadRequest)
		return
	}

	// Extract invoice codes from the data
	var invoiceCodes []string
	for _, item := range payload.Data {
		invoiceCodes = append(invoiceCodes, item.InvoiceCode)
	}

	// Check if all invoice codes exist in the database
	exists, err := repository.CheckInvoicesExist(invoiceCodes)

	if err != nil {
		helpers.SetResponse(w, r, "Database error", err.Error(), http.StatusInternalServerError)
		return
	}

	if !exists {
		helpers.SetResponse(w, r, "One or more invoice codes do not exist", nil, http.StatusBadRequest)
		return
	}

	// Transform SyncData into SyncDataUpdate
	var updates []models.SyncDataUpdate
	for _, item := range payload.Data {
		layout := "2006-01-02 15:04:05" // Expected format of time string
		attendTime, err := time.Parse(layout, item.AttendTime)
		if err != nil {
			helpers.SetResponse(w, r, "Invalid time format for attend_time", nil, http.StatusBadRequest)
			log.Println("Error parsing AttendTime:", err)
			return
		}
		updates = append(updates, models.SyncDataUpdate{
			InvoiceCode: item.InvoiceCode,
			AttendTime:  attendTime,
		})
	}

	// Update attendance status
	if err := repository.UpdateAttendanceStatus(updates, 1); err != nil {
		log.Println("Error", err)
		helpers.SetResponse(w, r, "Failed to update attendance status", nil, http.StatusInternalServerError)
		return
	}

	helpers.SetResponse(w, r, "Attendance status updated successfully", nil, http.StatusOK)
}
