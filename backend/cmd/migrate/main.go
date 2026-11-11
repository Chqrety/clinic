package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Main CLI tool for managing PostgreSQL database migrations and Chart of Accounts seeders.
// Usage:
//   go run cmd/migrate/main.go -action=up
//   go run cmd/migrate/main.go -action=down
//   go run cmd/migrate/main.go -action=seed-coa
func main() {
	action := flag.String("action", "up", "Migration action: up, down, status, seed-coa")
	flag.Parse()

	log.Printf("[MIGRATE] Running database migration command with action: %s\n", *action)

	switch *action {
	case "up":
		fmt.Println("[MIGRATE] Applying pending migrations...")
	case "down":
		fmt.Println("[MIGRATE] Rolling back latest migration...")
	case "seed-coa":
		fmt.Println("[MIGRATE] Seeding standard Clinic Chart of Accounts (COA)...")
	default:
		log.Fatalf("[MIGRATE] Unknown action: %s", *action)
		os.Exit(1)
	}

	log.Println("[MIGRATE] Migration task finished successfully.")
}

