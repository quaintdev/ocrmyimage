package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/otiai10/gosseract/v2"
)

func main() {
	// Serve static files
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/", fs)

	http.HandleFunc("/upload", uploadHandler)
	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("ocr request receieved")
	err := r.ParseMultipartForm(32 << 20) // Limit file size to 32MB
	if err != nil {
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	files := r.MultipartForm.File["images[]"]
	if len(files) == 0 {
		http.Error(w, "No files uploaded", http.StatusBadRequest)
		return
	}

	var results []string

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			http.Error(w, "Failed to open the file", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		// Create a temporary file to store the uploaded image
		tempFile, err := ioutil.TempFile("", "upload-*.png")
		if err != nil {
			http.Error(w, "Failed to create a temporary file", http.StatusInternalServerError)
			return
		}
		defer os.Remove(tempFile.Name())

		// Write the uploaded image to the temporary file
		_, err = io.Copy(tempFile, file)
		if err != nil {
			http.Error(w, "Failed to write the image to the temporary file", http.StatusInternalServerError)
			return
		}

		log.Println("file uploaded complete. running OCR...")

		client := gosseract.NewClient()
		client.SetLanguage(r.FormValue("language"))
		defer client.Close()
		client.SetImage(tempFile.Name())
		text, _ := client.Text()

		text = strings.ReplaceAll(text, "\n", " ")
		results = append(results, text)
	}

	// Create a struct to hold the OCR results
	type OCRResult struct {
		Texts []string `json:"texts"`
	}

	// Parse the OCR results into the struct
	ocrResult := OCRResult{Texts: results}
	log.Println(ocrResult)
	// Write the OCR results as JSON to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ocrResult)
	if err != nil {
		http.Error(w, "Failed to write the OCR results to the response", http.StatusInternalServerError)
		return
	}
}
