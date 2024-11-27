package api

import (
	"attendance-app/internal/helpers"
	"attendance-app/internal/models"
	"attendance-app/internal/repository"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

// LoginHandler handles operator login
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SetResponse(w, "Invalid request method", nil, http.StatusBadRequest)
		return
	}

	var operator models.Operator
	if err := json.NewDecoder(r.Body).Decode(&operator); err != nil {
		helpers.SetResponse(w, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Check if operator exists
	storedOperator, err := repository.GetOperatorByEmail(operator.Email)
	if err != nil {
		helpers.SetResponse(w, "Invalid email or password", err, http.StatusUnauthorized)
		return
	}

	log.Println(storedOperator)

	// Compare the stored hashed password with the input password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(operator.Password), bcrypt.DefaultCost)
	log.Printf("Hashed Password: %s", hashedPassword)

	err = bcrypt.CompareHashAndPassword([]byte(storedOperator.Password), []byte(operator.Password))
	if err != nil {
		log.Println("error")
		log.Println(err)
		helpers.SetResponse(w, "Invalid email or password", err, http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": storedOperator.Email,
		"id":    storedOperator.ID,
		"exp":   time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		helpers.SetResponse(w, "Could not create JWT token", nil, http.StatusInternalServerError)
		return
	}

	// Respond with the JWT token
	helpers.SetResponse(w, "Login successful", map[string]string{"token": tokenString}, http.StatusOK)
}

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
	if repository.OperatorExists(operator.Email, operator.Phone) {
		helpers.SetResponse(w, "The email or phone you entered was already registered", nil, http.StatusConflict)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(operator.Password), bcrypt.DefaultCost)
	log.Printf("Hashed Password: %s", hashedPassword)

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
