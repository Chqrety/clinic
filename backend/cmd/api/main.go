package main

import (
	"log"

	"clinic-api/internal/bootstrap"
)

func main() {
	// Native Direct Runner for Clinic Accounting & Financial ERP
	app, err := bootstrap.NewApplication()
	if err != nil {
		log.Fatalf("[FATAL] Failed to initialize application: %v", err)
	}

	if err := app.Run(); err != nil {
		log.Fatalf("[FATAL] Application terminated with error: %v", err)
	}
}
