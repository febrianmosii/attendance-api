package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

const (
	transactionsCount     = 10000 // Number of transactions
	detailsPerTransaction = 2     // Details per transaction
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	} else {
		log.Println("Successfully loaded .env file")
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	db_name := os.Getenv("DB_NAME")

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		username,
		password,
		host,
		port,
		db_name,
	))

	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("Starting data generation...")

	// Generate and insert transactions
	log.Println("Generating transactions...")
	if err := seedTransactions(conn, transactionsCount); err != nil {
		log.Fatalf("Error seeding transactions: %v", err)
	}

	// Generate and insert transaction details
	log.Println("Generating transaction details...")
	if err := seedTransactionDetails(conn, transactionsCount, detailsPerTransaction); err != nil {
		log.Fatalf("Error seeding transaction details: %v", err)
	}

	log.Println("Data generation completed!")
}

func seedTransactions(conn *pgx.Conn, count int) error {
	batchSize := 1000 // Insert in batches for better performance
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i += batchSize {
		end := i + batchSize
		if end > count {
			end = count
		}

		batch := &pgx.Batch{}
		for j := i; j < end; j++ {
			amount := rand.Float64()*1000 + 1                                // Random amount between 1 and 1000
			transactionDate := time.Now().AddDate(0, 0, -rand.Intn(30))      // Random date in the past 30 days
			invoiceCode := fmt.Sprintf("INV-%06d-%d", j+1, rand.Intn(10000)) // Add a random number to avoid duplicates
			status := []string{"pending", "paid", "failed"}[rand.Intn(3)]    // Random status
			userID := rand.Intn(2) + 1                                       // Random user ID

			batch.Queue(`
				INSERT INTO t_transactions (amount, transaction_date, invoice_code, transaction_status, user_id, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
			`, amount, transactionDate, invoiceCode, status, userID)
		}

		// Execute the batch
		results := conn.SendBatch(context.Background(), batch)
		if err := results.Close(); err != nil {
			return err
		}
	}

	return nil
}

func seedTransactionDetails(conn *pgx.Conn, transactionCount, detailsPerTransaction int) error {
	batchSize := 1000 // Insert in batches for better performance
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < transactionCount; i += batchSize {
		end := i + batchSize
		if end > transactionCount {
			end = transactionCount
		}

		batch := &pgx.Batch{}
		for j := i; j < end; j++ {
			for k := 0; k < detailsPerTransaction; k++ {
				transactionID := j + 1
				classID := rand.Intn(4) + 1
				ticketCode := fmt.Sprintf("TICKET-%06d-%d", transactionID, k) // Unique ticket code
				attendStatus := false
				attendTime := time.Time{} // NULL
				attendOperatorID := rand.Intn(2) + 1
				latestSyncAt := time.Time{} // NULL

				batch.Queue(`
					INSERT INTO t_transaction_details (transaction_id, class_id, ticket_code, attend_status, attend_time, attend_operator_id, latest_sync_at, created_at, updated_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
				`, transactionID, classID, ticketCode, attendStatus, attendTime, attendOperatorID, latestSyncAt)
			}
		}

		// Execute the batch
		results := conn.SendBatch(context.Background(), batch)
		if err := results.Close(); err != nil {
			return err
		}
	}

	return nil
}
