package main

import (
	"database/sql"
	"os"

	"log"

	"net/http"

	"time"

	database "example.com/tiernow/db/sqlc"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	_ "modernc.org/sqlite"
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
	dbConnection, err := sql.Open("sqlite", "./data.db")
	if err != nil {
		log.Fatalln(err)
	}

	    driver, err := sqlite.WithInstance(dbConnection, &sqlite.Config{})
    if err != nil {
        log.Fatal(err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://./db/migrations",
        "sqlite", driver,
    )
    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }

    log.Println("Migrations ran successfully")

	return database.New(dbConnection)
}

func initializeS3() *minio.Client {
	endpoint := "tiernow-garage:3900"
	accessKeyId := os.Getenv("GARAGE_KEY_ID")
	secretAccessKey := os.Getenv("GARAGE_SECRET_KEY")

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
