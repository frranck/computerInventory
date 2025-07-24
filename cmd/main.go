package main

import (
	"log"
	"os"

	"computerInventory/internal/adapter/db"
	"computerInventory/internal/adapter/rest"
	"computerInventory/internal/notifier"
	"computerInventory/internal/usecase"

	"github.com/gin-gonic/gin"
)

func main() {
	repo, err := db.NewSQLiteRepo("computers.db")
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	notifierURL := os.Getenv("NOTIFICATION_URL")
	if notifierURL == "" {
		notifierURL = "http://localhost:8080"
	}
	notifier := notifier.NewNotifier(notifierURL)
	service := usecase.NewService(repo, notifier)

	router := gin.Default()
	handler := rest.NewHandler(service)
	handler.RegisterRoutes(router)

	log.Println("Server running on :3000")
	if err := router.Run(":3000"); err != nil {
		log.Fatal(err)
	}
}
