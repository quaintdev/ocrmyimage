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
	log.Fatal(http.ListenAndServe(":8080", nil))
}
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Failed to retrieve the image", http.StatusBadRequest)
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

	log.Println("file uploaded complete. running ocr...")

	client := gosseract.NewClient()
	client.SetLanguage(r.FormValue("language"))
	defer client.Close()
	client.SetImage(tempFile.Name())
	text, _ := client.Text()

	text = strings.ReplaceAll(text, "\n", "")

	// Create a struct to hold the OCR output
	type OCRResult struct {
		Text string `json:"text"`
	}

	// Parse the JSON output from Tesseract
	var ocrResult OCRResult
	ocrResult.Text = text

	// Write the OCR result as JSON to the response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(ocrResult)
	if err != nil {
		http.Error(w, "Failed to write the OCR result to the response", http.StatusInternalServerError)
		return
	}
}
