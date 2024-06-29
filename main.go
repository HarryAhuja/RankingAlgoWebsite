package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	// HTTP handler to serve the image file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Replace this with your publicly accessible URL provided by S3
		objectURL := "https://rankingalgo.s3.eu-north-1.amazonaws.com/laptop.jpeg"

		// Fetch the object from S3
		resp, err := http.Get(objectURL)
		if err != nil {
			log.Printf("Error fetching object from S3: %v", err)
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Set Content-Type header based on the response from S3 or your file type
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

		// Copy object content directly to response writer
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error streaming file: %v", err)
			http.Error(w, "Failed to stream file", http.StatusInternalServerError)
		}
	})

	// Start HTTP server
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
