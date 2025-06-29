package main

import (
	"database/sql"
	"os"

	"log"

	"net/http"

	"time"

	database "example.com/tiernow/db/sqlc"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StorageInterface struct {
	db    *database.Queries
	minio *minio.Client
}

func initializeStorage() *StorageInterface {
	return &StorageInterface{
		db:    initializeDatabase(),
		minio: initializeS3(),
	}
}

func initializeDatabase() *database.Queries {
	dbConnection, err := sql.Open("sqlite3", "./data.db")
	if err != nil {
		log.Fatalln(err)
	}
	return database.New(dbConnection)
}

func initializeS3() *minio.Client {
	endpoint := "localhost:3900"
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Region: "garage",
	})

	if err != nil {
		log.Fatalln(err)
	}
	return minioClient
}

func main() {
	initializeS3()

	storage := initializeStorage()

	r := mux.NewRouter()
	r.HandleFunc("/images/{key}", storage.getImage).Methods("GET")
	r.HandleFunc("/tierlist", storage.createTierlist).Methods("POST")
	r.HandleFunc("/tierlist/{uuid}", storage.getTierlist).Methods("GET")
	r.HandleFunc("/tierlist/{uuid}/upload", storage.uploadImage).Methods("POST")
	r.HandleFunc("/tierlist/{uuid}/move", storage.moveImageToTier).Methods("POST")

	port := "5452"

	srv := &http.Server{
		Handler:      r,
		Addr:         ":" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(srv.ListenAndServe())
}
