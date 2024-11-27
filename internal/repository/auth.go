package repository

import (
	db "attendance-app/internal/database"
	"attendance-app/internal/models"
	"database/sql"
	"errors"
	"log"
	"time"
)

func InsertOperator(name, email, phone, hashedPassword string) (int, error) {
	query := `
		INSERT INTO m_operators (name, email, phone, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	var id int
	err := db.DB.QueryRow(query, name, email, phone, hashedPassword, time.Now(), time.Now()).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

func OperatorExists(email string, phone string) bool {
	var id int
	err := db.DB.QueryRow("SELECT id FROM m_operators WHERE email = $1 or phone = $2", email, phone).Scan(&id)

	return err == nil // If no error, operator exists
}

func GetOperatorByEmail(email string) (*models.Operator, error) {
	var operator models.Operator
	log.Printf("Searching for operator with email: %s", email)

	err := db.DB.QueryRow("SELECT id, name, email, phone, password FROM m_operators WHERE email = $1", email).Scan(&operator.ID, &operator.Name, &operator.Email, &operator.Phone, &operator.Password)
	if err != nil {
		log.Println("INI")
		if err == sql.ErrNoRows {
			log.Println("ITU")
			return nil, errors.New("operator not found")
		}
		log.Println("EG")
		return nil, err
	}
	log.Println("WKWK")
	return &operator, nil
}
