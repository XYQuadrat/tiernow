package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"

	schema "example.com/tiernow/db/sqlc"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/minio/minio-go/v7"
)

// 10MB
const MAX_UPLOAD_FILE_BYTES = 10 << 20

// Helper function to send JSON responses
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"failed to marshal JSON response"}`))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (storage *StorageInterface) uploadImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(MAX_UPLOAD_FILE_BYTES)
	file, handler, err := r.FormFile("image")

	if err != nil {
		http.Error(w, "Error when uploading: Couldn't retrieve file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageId := uuid.New().String()
	fileExt := strings.ToLower(strings.Split(handler.Filename, ".")[1])
	fileKey := fmt.Sprintf("%s.%s", imageId, fileExt)
	objectKey := fmt.Sprintf("images/%s", fileKey)

	uploadToS3(w, file, objectKey, handler.Header.Get("Content-Type"), storage.minio)

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	entry, err := storage.db.UploadImageMetadata(r.Context(), schema.UploadImageMetadataParams{TierlistUuid: uuid, FileKey: fileKey})

	if err != nil {
		log.Printf("Failed to create file metadata: %v", err)
		http.Error(w, "Error when uploading: Couldn't save metadata", http.StatusInternalServerError)
		return
	}

	response := struct{
		ID	int64 `json:"id"`
		Filename string `json:"filename"`
	}{
		ID: entry.ID,
		Filename: fileKey,
	}

	respondWithJSON(w, http.StatusCreated, response)
}

func uploadToS3(w http.ResponseWriter, file multipart.File, key string, contentType string, minioClient *minio.Client) {
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
}

func (storage *StorageInterface) getImage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := "images/" + vars["key"]

	image, err := storage.minio.GetObject(context.Background(), "tiernow", key, minio.GetObjectOptions{})

	if err != nil {
		log.Printf("Failed to retrieve file: %v", err)
		http.Error(w, "Error when retrieving: Couldn't find image", http.StatusNotFound)
		return
	}

	metadata, err := storage.minio.StatObject(context.Background(), "tiernow", key, minio.StatObjectOptions{})
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


func (storage *StorageInterface) getTierlist(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]

	tierlist, err := storage.db.GetTierlist(r.Context(), uuid)

	if r, ok := tierlist.Tiers.(string); ok {
		jsonMessage := json.RawMessage(r)
		tierlist.Tiers = jsonMessage
	}
	if r, ok := tierlist.UnassignedEntries.(string); ok {
		jsonMessage := json.RawMessage(r)
		tierlist.UnassignedEntries = jsonMessage
	}

	if err != nil {
		http.Error(w, "Couldn't get tierlist", http.StatusNotFound)
		return
	}
	respondWithJSON(w, http.StatusOK, tierlist)
}

func (storage *StorageInterface) createTierlist(w http.ResponseWriter, r *http.Request) {
	var requestBody schema.CreateTierlistParams
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := schema.CreateTierlistParams{
		Uuid: uuid.NewString(),
		Name: requestBody.Name,
	}

	tierlist, err := storage.db.CreateTierlist(r.Context(), params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defaultTiers := []string{"S", "A", "B", "C", "D"}

	for i := range defaultTiers {
		log.Default().Print("Creating default tier")
		tierParams := schema.CreateTierParams{
			TierlistUuid: params.Uuid,
			Name:         defaultTiers[i],
			Order:        int64(i),
		}
		_, err := storage.db.CreateTier(r.Context(), tierParams)

		if err != nil {
			http.Error(w, "Error: Could not create default tier", http.StatusBadRequest)
			return
		}
	}

	respondWithJSON(w, http.StatusCreated, tierlist)
}

func (storage *StorageInterface) moveImageToTier(w http.ResponseWriter, r *http.Request) {
	var requestBody schema.SetImageTierParams
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	params := schema.SetImageTierParams{
		TierID: requestBody.TierID,
		ID: requestBody.ID,
	}

	if params.TierID != nil {
		log.Default().Printf("Moved image with ID %d to tier %d",params.ID, *params.TierID)
	} else {
		log.Default().Printf("Moved image with ID %d to Uploaded Items",params.ID)
	}

	entry, err := storage.db.SetImageTier(r.Context(), params)
	if err != nil {
		http.Error(w, "Error: Could not move image to tier", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, entry)
}