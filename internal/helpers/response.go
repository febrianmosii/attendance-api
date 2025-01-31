package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"unicode"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

// Response struct to structure our API responses
type Response struct {
	ResponseID string      `json:"response_id"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"`
	Errors     interface{} `json:"errors,omitempty"`
}

// Set up a logger
var responseLogger *log.Logger

func init() {
	// Ensure the log folder exists
	logFolder := "logs"
	if err := os.MkdirAll(logFolder, os.ModePerm); err != nil {
		log.Fatalf("Failed to create log folder: %v", err)
	}

	// Generate the file name with the current date in -YYMMDD format
	currentDate := time.Now().Format("060102") // YYMMDD format
	logFileName := logFolder + "/responses-" + currentDate + ".log"

	// Create or open the log file
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	// Initialize the logger
	responseLogger = log.New(logFile, "", log.LstdFlags)
}

// SetResponse formats and sends JSON responses while logging request and response
func SetResponse(w http.ResponseWriter, r *http.Request, message string, data interface{}, httpCode int) {
	body := ""

	// Log the request
	requestLog := map[string]interface{}{
		"body":    string(body),
		"headers": r.Header,
	}

	// Generate a unique response ID
	responseID := generateUUID()

	// Create the response struct
	response := Response{
		ResponseID: responseID,
		Message:    message,
	}

	// Check HTTP status code and set data/errors
	if httpCode >= 200 && httpCode < 300 {
		if data != nil {
			response.Data = data
		}
	} else {
		if httpCode == 422 {
			validationResponse := FormatValidationErrors(data)

			if data != nil {
				response.Errors = validationResponse
			}
		} else {
			if data != nil {
				response.Errors = data
			}
		}
	}

	// Combine request and response for logging
	logEntry := map[string]interface{}{
		"request": requestLog,
		"response": map[string]interface{}{
			"response_id": responseID,
			"message":     message,
		},
		"http_code": httpCode,
	}

	if data != nil {
		logEntry["response"].(map[string]interface{})["data"] = data
	}

	// Set response headers and write the JSON response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
	}
}

func generateUUID() string {
	return uuid.New().String()
}

func FormatValidationErrors(err interface{}) map[string][]string {
	validationErrors := make(map[string][]string)

	for _, fieldErr := range err.(validator.ValidationErrors) {
		fieldName := ToSnakeCase(fieldErr.Field())
		message := ""
		switch fieldErr.Tag() {
		case "required":
			message = fieldName + " is required."
		case "max":
			message = fieldName + " cannot exceed " + fieldErr.Param() + " characters."
		}
		validationErrors[fieldName] = append(validationErrors[fieldName], message)
	}

	return validationErrors
}

func ToSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}
