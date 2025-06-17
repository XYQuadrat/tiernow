package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// 10MB
const MAX_UPLOAD_FILE_BYTES = 10 << 20

var (
	minioClient *minio.Client
)

func initializeS3() {
	endpoint := "localhost:3900"
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")

	var err error
	minioClient, err = minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Region: "garage",
	})

	if err != nil {
		log.Fatalln(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(MAX_UPLOAD_FILE_BYTES)
	file, handler, err := r.FormFile("image")

	if err != nil {
		http.Error(w, "Error when uploading: Couldn't retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageId := uuid.New().String()
	fileExt := strings.ToLower(strings.Split(handler.Filename, ".")[1])
	objectKey := fmt.Sprintf("images/%s.%s", imageId, fileExt)

	uploadToS3(w, file, objectKey, handler.Header.Get("Content-Type"))
}

func uploadToS3(w http.ResponseWriter, file multipart.File, key string, contentType string) {
	_, err := minioClient.PutObject(
		context.Background(),
		"tiernow",
		key,
		file,
		-1,
		minio.PutObjectOptions{ContentType: contentType},
	)

	if err != nil {
		log.Printf("Failed to upload file: %v", err)
		http.Error(w, "Error when uploading: Couldn't upload to storage", http.StatusInternalServerError)
		return
	}

	log.Printf("Uploaded new file to %v", key)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(key))
}

func retrieveHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := "images/" + vars["key"]

	image, err := minioClient.GetObject(context.Background(), "tiernow", key, minio.GetObjectOptions{})

	if err != nil {
		log.Printf("Failed to retrieve file: %v", err)
		http.Error(w, "Error when retrieving: Couldn't find image", http.StatusNotFound)
		return
	}

	metadata, err := minioClient.StatObject(context.Background(), "tiernow", key, minio.StatObjectOptions{})
	if err != nil {
		log.Printf("Failed to retrieve file metadata: %v", err)
		http.Error(w, "Error when retrieving: Couldn't get file metadata", http.StatusNotFound)
		return
	}
	defer image.Close()

	log.Printf("Retrieved file %v with content type: %v", key, metadata.ContentType)
	w.Header().Set("Content-Type", metadata.ContentType)
	io.Copy(w, image)
}

func main() {
	initializeS3()

	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/images/{key}", retrieveHandler).Methods("GET")

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
