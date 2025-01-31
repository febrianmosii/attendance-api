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

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey = []byte(os.Getenv("JWT_SECRET"))

var validate = validator.New()

type LoginRequest struct {
	UserName string  `json:"username" validate:"required,max=255"`
	Password string  `json:"password" validate:"required,max=255"`
	DeviceId *string `json:"device_id" validate:"required"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		helpers.SetResponse(w, r, "Invalid request method", nil, http.StatusBadRequest)
		return
	}

	// Decode the JSON request
	var loginReq LoginRequest

	// Decode JSON and validate fields
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		helpers.SetResponse(w, r, "Invalid request body", nil, http.StatusBadRequest)
		return
	}

	// Validate the request body
	validate := validator.New()

	err := validate.Struct(loginReq)
	if err != nil {
		helpers.SetResponse(w, r, "Validation Failed", err, http.StatusUnprocessableEntity)
		return
	}

	// Check if operator exists
	storedOperator, err := repository.GetOperatorByUsername(loginReq.UserName)

	if err != nil {
		log.Printf("Authentication error: %v", err)
		helpers.SetResponse(w, r, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	// Compare device id to prevent multi login user
	print("storedOperator.DeviceId", loginReq.DeviceId)

	if *loginReq.DeviceId != *storedOperator.DeviceId {
		helpers.SetResponse(w, r, "You have logged in from another device. Please contact the administrator for further assistance.", nil, http.StatusForbidden)
		return
	}

	// Compare the stored hashed password with the input password
	err = bcrypt.CompareHashAndPassword([]byte(storedOperator.Password), []byte(loginReq.Password))

	if err != nil {
		log.Printf("Password mismatch for username: %s", loginReq.UserName)
		helpers.SetResponse(w, r, "Invalid username or password", nil, http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": storedOperator.UserName,
		"id":       storedOperator.ID,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	})

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("JWT_SECRET_KEY environment variable is not set")
	}

	tokenString, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		log.Printf("JWT signing error: %v", err)
		helpers.SetResponse(w, r, "Could not create JWT token", nil, http.StatusInternalServerError)
		return
	}

	storedOperator.DeviceId = loginReq.DeviceId
	storedOperator.DeviceAccessToken = &tokenString

	// Respond with the JWT token
	response := models.LoginResponseData{
		Token: tokenString,
		User:  *storedOperator,
	}

	if err := repository.UpdateOperatorDeviceInformation(storedOperator.DeviceId, storedOperator.DeviceAccessToken, storedOperator.ID); err != nil {
		log.Println("Error", err)
		helpers.SetResponse(w, r, "Failed to update device information", nil, http.StatusInternalServerError)
		return
	}

	helpers.SetResponse(w, r, "Login Success!", response, http.StatusOK)
}

// isValidEmail validates the email format
func isValidEmail(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(regex)
	return re.MatchString(email)
}
