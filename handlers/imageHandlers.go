package handlers

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
)

type ImageHandler struct {
	Ctx context.Context
	DB  *sql.DB
}

// ServeHTTP handles requests for serving images.
func (i ImageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	imageName := r.URL.Path
	// Query MySQL to get the URL for imageName
	query := "SELECT url FROM image_table WHERE name = ?"
	var imageURL string
	err := i.DB.QueryRowContext(i.Ctx, query, imageName).Scan(&imageURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Fetch the image from imageURL
	resp, err := http.Get(imageURL)
	if err != nil {
		log.Printf("Error fetching image from URL: %v", err)
		http.Error(w, "Failed to retrieve image", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set Content-Type header based on the response from imageURL or your file type
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))

	// Copy image content directly to response writer
	if _, err := io.Copy(w, resp.Body); err != nil {
		log.Printf("Error streaming image: %v", err)
		http.Error(w, "Failed to stream image", http.StatusInternalServerError)
	}
}
