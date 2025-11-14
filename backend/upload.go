package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func UploadImageHandler(w http.ResponseWriter, r *http.Request) {
	// Limit request size (optional but recommended)
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20) // 10MB max

	// Parse multipart form
	err := r.ParseMultipartForm(20 << 20) // 20MB
	if err != nil {
		http.Error(w, "Error parsing form: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Retrieve the file (the key must match: formData.append("imageFile", file))
	file, fileHeader, err := r.FormFile("imageFile")
	if err != nil {
		http.Error(w, "Could not get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// OPTIONAL: generate a random file name
	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	// Make sure uploads folder exists
	os.MkdirAll("./uploads", 0755)

	// Create the destination file
	dst, err := os.Create("./uploads/" + filename)
	if err != nil {
		http.Error(w, "Could not save file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file to destination
	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Could not write file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Response (JSON)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"status":"success", "file":"%s"}`, filename)
}
