package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// Response struct to structure our API responses
type Response struct {
	ResponseID string      `json:"response_id"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

// SetResponse formats and sends JSON responses with optional logging
func SetResponse(w http.ResponseWriter, message string, data interface{}, httpCode int) {
	// Generate a unique response ID (you can use a UUID generator here)
	responseID := generateUUID()

	// Create the response struct
	response := Response{
		ResponseID: responseID,
		Message:    message,
	}

	// Check if data should be included or errors should be included based on HTTP status code
	if httpCode >= 200 && httpCode < 300 {
		// Success case: Attach data to response
		if data != nil {
			response.Data = data
		}
	} else if httpCode == 422 {
		// Validation error case: Attach errors to response
		if data != nil {
			response.Errors = data
		}
	} else {
		// For non-success HTTP codes (500, etc.): Attach error details
		if data != nil {
			response.Errors = data
		}
	}

	// Set the response header to application/json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)

	// Send the JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
	}
}

// Generate a simple UUID-like response ID (you can use a proper UUID package)
func generateUUID() string {
	return uuid.New().String()
}
