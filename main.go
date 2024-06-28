package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Config struct {
	AWSAccessKeyID     string `json:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `json:"AWS_SECRET_ACCESS_KEY"`
	AWSRegion          string `json:"AWS_REGION"`
	S3Bucket           string `json:"S3_BUCKET"`
}

func loadConfig(filename string) (Config, error) {
	var config Config
	configFile, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&config)
	return config, err
}

func main() {
	// Load AWS credentials and region from config file
	config, err := loadConfig("config.json")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Create AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AWSRegion),
		Credentials: credentials.NewStaticCredentials(
			config.AWSAccessKeyID,
			config.AWSSecretAccessKey,
			""),
	})
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	svc := s3.New(sess)

	// HTTP handler to serve the image file
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileKey := "laptop.jpeg"

		// Get the object from S3
		resp, err := svc.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(config.S3Bucket),
			Key:    aws.String(fileKey),
		})
		if err != nil {
			log.Printf("Error getting object %s from bucket rankingalgo: %v", fileKey, err)
			http.Error(w, "Failed to retrieve file", http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Set Content-Type header based on S3 response or your file type
		w.Header().Set("Content-Type", aws.StringValue(resp.ContentType))

		// Copy object content directly to response writer
		if _, err := io.Copy(w, resp.Body); err != nil {
			log.Printf("Error streaming file: %v", err)
			http.Error(w, "Failed to stream file", http.StatusInternalServerError)
		}
	})

	// Start HTTP server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
