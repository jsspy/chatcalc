package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "POST only", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // 50MB
	if err != nil {
		http.Error(w, "parse error", 500)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file error", 500)
		return
	}
	defer file.Close()

	os.MkdirAll("uploads", 0755)
	filename := fmt.Sprintf("%d_%s", os.Getpid(), header.Filename)
	filepath := filepath.Join("uploads", filename)

	out, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "save error", 500)
		return
	}
	defer out.Close()

	io.Copy(out, file)

	url := "/uploads/" + filename

	resp := map[string]string{"url": url}
	json.NewEncoder(w).Encode(resp)
}
