package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var connStr string
var appName string

func init() {
	host := os.Getenv("MONGO_HOST")
	username := os.Getenv("MONGO_USERNAME")
	password := os.Getenv("MONGO_PASSWORD")
	appName = os.Getenv("MONGO_APPNAME")
	if host == "" || username == "" || password == "" || appName == "" {
		panic("Missing required environment variables")
	}

	connStr = fmt.Sprintf("mongodb+srv://%s:%s@%s/?appname=%s", username, password, host, appName)
}

func main() {

	// Connect to the MongoDB server
	client, err := mongo.Connect(options.Client().ApplyURI(connStr))
	if err != nil {
		panic(fmt.Errorf("Failed to connect to MongoDB: %w", err))
	}
	defer client.Disconnect(nil)

	// Ping the server to verify the connection
	if err = client.Ping(nil, nil); err != nil {
		panic(fmt.Errorf("Failed to ping MongoDB: %w", err))
	}

	// Use the proper database
	db := client.Database(appName)

	// Create a GridFS bucket if it doesn't exist
	bucket := db.GridFSBucket()

	// Open the file we want to upload via a stream
	file, err := os.Open("main.go")
	if err != nil {
		panic(fmt.Errorf("Failed to open file: %w", err))
	}
	defer file.Close()

	// Upload the file to GridFS via a stream.
	// Parts of the file will still be read into memory, but this is
	// much more efficient than reading the entire file into memory at once.
	objectID, err := bucket.UploadFromStream(context.Background(), "main.go", file, nil)
	if err != nil {
		panic(err)
	}
	slog.Info("New file uploaded", "name", "main.go", "objectID", objectID.Hex())

	// Download the file from GrifFS with the object ID appended to the filename
	fname := fmt.Sprintf("main_%s.go.dl", objectID.Hex())
	downloadedFile, err := os.Create(fname)
	if err != nil {
		panic(fmt.Errorf("Failed to create file: %w", err))
	}
	defer downloadedFile.Close()
	_, err = bucket.DownloadToStream(context.Background(), objectID, downloadedFile)
	if err != nil {
		panic(fmt.Errorf("Failed to download file: %w", err))
	}
	slog.Info("File downloaded", "name", fname, "objectID", objectID.Hex())
	// Delete the file from GridFS
	err = bucket.Delete(context.Background(), objectID)
	if err != nil {
		panic(fmt.Errorf("Failed to delete file: %w", err))
	}
	slog.Warn("File deleted from gridfs", "name", "main.go", "objectID", objectID.Hex())
}
