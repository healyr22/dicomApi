package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nfnt/resize"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

var storagePath = "./uploads"

// Upload and store the DICOM file
func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the uploaded file
	err := r.ParseMultipartForm(10 << 20) // Limit upload size to 10 MB
	if err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Save the file
	filePath := filepath.Join(storagePath, header.Filename)
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()
	io.Copy(out, file)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "File uploaded and stored successfully"})
}

// Extract a DICOM header attribute
func handleExtractHeader(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")
	tagName := r.URL.Query().Get("tag")

	if fileName == "" || tagName == "" {
		http.Error(w, "Both filePath and tag are required", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(storagePath, fileName)

	// Load the DICOM file
	dataset, err := dicom.ParseFile(filePath, nil)
	if err != nil {
		http.Error(w, "Failed to parse DICOM file", http.StatusInternalServerError)
		return
	}

	// Find the tag value using FindByName
	tagInfo, err := tag.FindByName(tagName)
	if err != nil {
		http.Error(w, "Invalid DICOM tag", http.StatusBadRequest)
		return
	}

	// Find the tag element by its Tag
	element, err := dataset.FindElementByTag(tagInfo.Tag)
	if err != nil {
		http.Error(w, "Tag not found in DICOM file", http.StatusNotFound)
		return
	}

	// Respond with the tag value
	json.NewEncoder(w).Encode(map[string]interface{}{
		"tag":   tagName,
		"value": element.String(),
	})
}

// Convert a DICOM file to PNG
func handleConvertToPNG(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Query().Get("fileName")

	if fileName == "" {
		http.Error(w, "filePath is required", http.StatusBadRequest)
		return
	}
	filePath := filepath.Join(storagePath, fileName)

	// Load the DICOM file
	dataset, err := dicom.ParseFile(filePath, nil)
	if err != nil {
		http.Error(w, "Failed to parse DICOM file", http.StatusInternalServerError)
		return
	}

	// Extract pixel data
	pixelDataElement, err := dataset.FindElementByTag(tag.PixelData)
	if err != nil {
		http.Error(w, "Failed to extract pixel data", http.StatusInternalServerError)
		return
	}

	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for i, frame := range pixelDataInfo.Frames {
		img, err := frame.GetImage()
		if err != nil {
			http.Error(w, "Failed to convert frame to image", http.StatusInternalServerError)
			return
		}

		// Resize for browser display
		resizedImg := resize.Resize(512, 512, img, resize.Lanczos3)

		// Write the first frame as PNG
		if i == 0 {
			w.Header().Set("Content-Type", "image/png")
			jpeg.Encode(w, resizedImg, &jpeg.Options{Quality: 100})
			return
		}
	}
}

func main() {
	// Ensure storage path exists
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		os.Mkdir(storagePath, os.ModePerm)
	}

	// Register handlers
	http.HandleFunc("/upload", handleUpload)
	http.HandleFunc("/extract-header", handleExtractHeader)
	http.HandleFunc("/convert-to-png", handleConvertToPNG)

	// Start the server
	port := "8080"
	fmt.Printf("Server is running on http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
	}
}
