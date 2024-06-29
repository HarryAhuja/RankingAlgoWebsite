package main

import (
	"context"
	"log"
	"net/http"
	"rankingAlgoWebsite/db"
	"rankingAlgoWebsite/handlers"
)

func main() {

	db := db.ConfigureDB()
	if db == nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	imageHandler := handlers.ImageHandler{
		Ctx: context.Background(),
		DB:  db,
	}

	// Register handler for /get-product/ path
	http.Handle("/get-product/", http.StripPrefix("/get-product/", imageHandler))

	// Start HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
