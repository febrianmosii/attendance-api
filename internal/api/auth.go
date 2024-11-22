package api

import (
	"attendance-app/internal/helpers"
	"attendance-app/internal/models"
	"attendance-app/internal/repository"
	"encoding/json"
	"net/http"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// RegisterHandler handles operator registration
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SetResponse(w, "Invalid request method", nil, http.StatusBadRequest)
		return
	}

	var operator models.Operator
	if err := json.NewDecoder(r.Body).Decode(&operator); err != nil {
		helpers.SetResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Validate the form data
	errors := validateOperator(operator)
	if errors != nil {
		helpers.SetResponse(w, "Validation failed", errors, http.StatusOK)
		return
	}

	// Check if operator exists
	if repository.OperatorExists(operator.Email) {
		helpers.SetResponse(w, "Operator already exists with this email", nil, http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(operator.Password), bcrypt.DefaultCost)
	if err != nil {
		helpers.SetResponse(w, "Failed to hash password", nil, http.StatusInternalServerError)
		return
	}

	// Insert the operator into the database
	operatorID, err := repository.InsertOperator(operator.Name, operator.Email, operator.Phone, string(hashedPassword))
	if err != nil {
		helpers.SetResponse(w, "Failed to register operator", nil, http.StatusInternalServerError)
		return
	}

	// Respond with success
	operator.ID = operatorID
	operator.Password = "" // Clear password for the response
	helpers.SetResponse(w, "Operator registered successfully", operator, http.StatusCreated)
}

// validateOperator validates operator form data
func validateOperator(operator models.Operator) *models.RegisterValidationErrors {
	errors := &models.RegisterValidationErrors{}

	if operator.Name == "" {
		errors.NameError = "Name is required"
	}

	if operator.Email == "" {
		errors.EmailError = "Email is required"
	} else if !isValidEmail(operator.Email) {
		errors.EmailError = "Invalid email format"
	}

	if operator.Phone == "" {
		errors.PhoneError = "Phone number is required"
	}

	if operator.Password == "" {
		errors.PasswordError = "Password is required"
	} else if len(operator.Password) < 6 || len(operator.Password) > 12 {
		errors.PasswordError = "Password must be between 6 and 12 characters"
	}

	if operator.PasswordConfirmation == "" {
		errors.PasswordConfirmationError = "Password confirmation is required"
	} else if operator.Password != operator.PasswordConfirmation {
		errors.PasswordConfirmationError = "Passwords do not match"
	}

	if errors.NameError != "" || errors.EmailError != "" || errors.PhoneError != "" || errors.PasswordError != "" || errors.PasswordConfirmationError != "" {
		return errors
	}

	return nil
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
