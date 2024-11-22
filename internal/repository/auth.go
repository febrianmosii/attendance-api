package repository

import (
	db "attendance-app/internal/database"
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

func OperatorExists(email string) bool {
	var id int
	err := db.DB.QueryRow("SELECT id FROM m_operators WHERE email = $1", email).Scan(&id)
	return err == nil // If no error, operator exists
}
