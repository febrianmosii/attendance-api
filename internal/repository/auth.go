package repository

import (
	db "attendance-app/internal/database"
	"attendance-app/internal/models"
	"database/sql"
	"encoding/json"
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

func UpdateOperatorDeviceInformation(device_id *string, device_token *string, operator_id int) error {
	// Prepare the base query
	query := `
		UPDATE m_admin_attendances 
		SET 
			device_id = $1,
			device_access_token = $2,
			updated_at = now()
		WHERE id = $3`

	_, err := db.DB.Exec(query, device_id, device_token, operator_id)

	if err != nil {
		return err
	}

	return nil
}

func OperatorExists(email string, phone string) bool {
	var id int
	err := db.DB.QueryRow("SELECT id FROM m_operators WHERE email = $1 or phone = $2", email, phone).Scan(&id)

	return err == nil // If no error, operator exists
}

func GetOperatorByUsername(username string) (*models.Operator, error) {
	var result models.Operator

	// Main query for the operator and role details
	query := `
		SELECT 
			a.id, 
			a.name, 
			a.username, 
			a.email, 
			a.phone, 
			a.is_active, 
			a.is_limited_event_access, 
			a.is_limited_classes_access, 
			a.device_id, 
			a.device_access_token, 
			a.password,
			b.id as role_id,
			b.name as role_name
		FROM m_admin_attendances a
		LEFT JOIN m_admin_attendance_roles b ON b.id = a.role_id
		WHERE a.username = $1`

	// Fetch the main operator details
	err := db.DB.QueryRow(query, username).
		Scan(
			&result.ID,
			&result.Name,
			&result.UserName,
			&result.Email,
			&result.Phone,
			&result.IsActive,
			&result.IsLimitedEventAccess,
			&result.IsLimitedClassAccess,
			&result.DeviceId,
			&result.DeviceAccessToken,
			&result.Password,
			&result.Role.ID,
			&result.Role.Name,
		)
	if err != nil {
		return nil, err // Handle error appropriately
	}

	if result.IsLimitedEventAccess {
		// Eager load all events and classes related to the operator
		eagerQuery := `
			SELECT 
				a.event_id,
				d.event_name,
				b.class_id,
				c.class_name
			FROM t_list_show_events_admin a
			LEFT JOIN (
				SELECT 
					x.event_id, 
					t_list_show_classes_admin.class_id 
				FROM t_list_show_classes_admin
				JOIN m_classes x ON x.id = t_list_show_classes_admin.class_id
				GROUP by event_id, class_id
			) b ON a.event_id = b.event_id
			LEFT JOIN m_classes c ON c.id = b.class_id
			LEFT JOIN m_events d ON d.id = a.event_id
			WHERE a.admin_attendance_id = $1 and a.deleted_at is null`

		rows, err := db.DB.Query(eagerQuery, result.ID)

		if err != nil {
			return nil, err // Handle error appropriately
		}

		defer rows.Close()

		// Maps to group events and their associated classes
		eventsMap := make(map[int]*models.AccessEvent)
		classesMap := make(map[int][]models.EventClass)

		// Process the results
		for rows.Next() {
			var eventID, classID sql.NullInt64
			var eventName string
			var className *string

			// Scan each row
			err := rows.Scan(&eventID, &eventName, &classID, &className)

			if err != nil {
				return nil, err
			}

			// Add event to the map if not already present
			if eventID.Valid {
				if _, exists := eventsMap[int(eventID.Int64)]; !exists {
					eventsMap[int(eventID.Int64)] = &models.AccessEvent{
						ID:      int(eventID.Int64),
						Name:    eventName,
						Classes: []models.EventClass{},
					}
				}
			}

			// Add class to the classes map
			if classID.Valid {
				classesMap[int(eventID.Int64)] = append(classesMap[int(eventID.Int64)], models.EventClass{
					ID:   int(classID.Int64),
					Name: className,
				})
			}
		}

		// Combine the events and classes
		var accessEvents []models.AccessEvent
		for eventID, event := range eventsMap {
			if classes, exists := classesMap[eventID]; exists {
				event.Classes = classes
			} else {
				event.Classes = []models.EventClass{}
			}

			accessEvents = append(accessEvents, *event)
		}

		// Assign the result to the operator
		result.AccessEvents = accessEvents
	}

	log.Println(result)

	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Printf("Error marshalling to JSON: %v", err)
	}

	log.Println(string(jsonResult))

	return &result, nil
}
