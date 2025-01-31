package repository

import (
	db "attendance-app/internal/database"
	"attendance-app/internal/models"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

// GetAttendanceData fetches attendance data for a given operator ID.
func GetAttendanceData(operatorId int) ([]models.SyncResponse, error) {
	exists, err := IsOperatorExists(operatorId)

	if err != nil {
		log.Println("Error checking operator existence:", err)
		return nil, err
	}

	if !exists {
		log.Println("Operator not found with ID:", operatorId)
		return nil, errors.New("operator_not_found")
	}

	// Query to fetch attendance data based on operatorId
	query := `
        SELECT 
            a.id AS id, 
            b.invoice_code, 
            a.ticket_code as ticket_code,
            a.attend_status as attend_status, 
            a.attend_time as attend_time,
            g.id as event_id,
            g.name as event_name,
			f.id as class_id,
            f.name as class_name,
            e.id as user_id,
            e.name as user_name,
            e.phone as user_phone,
            e.email as user_email
        FROM 
            t_transaction_details a 
            LEFT JOIN t_transactions b ON b.id = a.transaction_id
            LEFT JOIN t_operators_classes c ON c.class_id = a.class_id
            LEFT JOIN m_operators d ON d.id = c.operator_id
            LEFT JOIN m_users e ON e.id = b.user_id
            LEFT JOIN m_classes f ON f.id = a.class_id
            LEFT JOIN m_events g ON g.id = f.event_id
        WHERE c.operator_id = $1
        ORDER BY 
            c.operator_id ASC;
	`

	rows, err := db.DB.Query(query, operatorId)
	if err != nil {
		log.Println("Error querying database:", err)
		return nil, err
	}
	defer rows.Close()

	var responses []models.SyncResponse

	for rows.Next() {
		var response models.SyncResponse
		var class models.Class
		var user models.User
		var event models.Event

		// Scan database fields into structs
		if err := rows.Scan(
			&response.ID,
			&response.InvoiceCode,
			&response.TicketCode,
			&response.AttendStatus,
			&response.AttendTime,
			&event.ID,
			&event.Name,
			&class.ID,
			&class.Name,
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Phone,
		); err != nil {
			log.Println("Error scanning row:", err)
			return nil, err
		}

		// Link Class and Event
		class.Event = event
		response.Class = class
		response.User = user

		responses = append(responses, response)
	}

	if err := rows.Err(); err != nil {
		log.Println("Error with rows:", err)
		return nil, err
	}

	return responses, nil
}

func IsOperatorExists(operatorId int) (bool, error) {
	query := `SELECT id from m_operators where id = $1`

	var id string

	err := db.DB.QueryRow(query, operatorId).Scan(&id)

	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return false, nil
		}

		// For other errors, log and return
		log.Println("Error querying operator existence:", err)
		return false, err
	}

	return true, nil
}

// CheckInvoicesExist checks if all invoice codes exist in the database.
func CheckInvoicesExist(invoiceCodes []string) (bool, error) {
	// Generate PostgreSQL placeholders dynamically
	placeholders := make([]string, len(invoiceCodes))
	args := make([]interface{}, len(invoiceCodes))

	for i, code := range invoiceCodes {
		placeholders[i] = "$" + strconv.Itoa(i+1) // PostgreSQL uses $1, $2, ...
		args[i] = code
	}

	query := `
		SELECT COUNT(*) 
		FROM t_transaction_details 
		WHERE ticket_code IN (` + strings.Join(placeholders, ",") + `)`

	log.Println("query", query)
	log.Println("args", args)

	var count int
	err := db.DB.QueryRow(query, args...).Scan(&count)
	if err != nil {
		log.Println("error", err)
		return false, err
	}

	return count == len(invoiceCodes), nil
}

func UpdateAttendanceStatus(data []models.SyncDataUpdate, operatorID int) error {
	// Prepare the base query and arguments slice
	query := `
		UPDATE t_transaction_details t
		SET 
			attend_status = true, 
			attend_time = v.attend_time,
			attend_operator_id = v.attend_operator_id,
			latest_sync_at = now()
		FROM (VALUES`

	args := []interface{}{}
	// Loop through the data to build the query and append the values
	for i, item := range data {
		// Add placeholder values for attend_time, operator_id, and invoice_code
		if i > 0 {
			query += ","
		}
		// Append values for this row
		query += fmt.Sprintf("($%d::timestamp, $%d::int, $%d)", i*3+1, i*3+2, i*3+3)

		// Append the corresponding arguments for each row
		args = append(args, item.AttendTime, operatorID, item.InvoiceCode)
	}

	// Complete the query by closing the FROM clause and adding the WHERE clause
	query += `) AS v(attend_time, attend_operator_id, ticket_code)
	WHERE t.ticket_code = v.ticket_code`

	// Execute the query
	_, err := db.DB.Exec(query, args...)
	return err
}

func UpdateAttendanceStatusOld(data []models.SyncDataUpdate, operatorID int) error {
	// Prepare the base query
	query := `
		UPDATE t_transaction_details
		SET 
			attend_status = true, 
			attend_time = $1,
			attend_operator_id = $2,
			latest_sync_at = now()
		WHERE ticket_code = $3`

	// Loop through the data and update each record individually
	for _, item := range data {
		// Execute the update query for each record
		_, err := db.DB.Exec(query, item.AttendTime, operatorID, item.InvoiceCode)
		if err != nil {
			// Log the error and continue, or return depending on the error handling strategy
			log.Printf("Error updating attendance status for ticket %s: %v", item.InvoiceCode, err)
			// Optionally, return the error if you want to stop on failure
			return err
		}
	}

	return nil
}
