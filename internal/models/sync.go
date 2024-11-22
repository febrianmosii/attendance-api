package models

import (
	"time"
)

type Event struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type Class struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Event Event  `json:"event"`
}

type SyncResponse struct {
	ID           int        `json:"id"`
	InvoiceCode  string     `json:"invoice_code"`
	TicketCode   string     `json:"ticket_code"`
	AttendStatus bool       `json:"attend_status"`
	AttendTime   *time.Time `json:"attend_time"`
	User         User       `json:"user"`
	Class        Class      `json:"class"`
}

type SyncData struct {
	InvoiceCode string `json:"invoice_code"`
	AttendTime  string `json:"attend_time"`
}

type SyncDataUpdate struct {
	InvoiceCode string    `json:"invoice_code"`
	AttendTime  time.Time `json:"attend_time"`
}

type SyncPayload struct {
	Data []SyncData `json:"data"`
}
