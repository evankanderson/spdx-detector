package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/CycloneDX/license-scanner/api/scanner"
)

func scanLicense(file []byte) ([]string, error) {
	scan := scanner.ScanSpecs{
		Specs: []scanner.ScanSpec{
			{
				Name:        "online",
				LicenseText: string(file),
			},
		},
	}

	results, err := scan.ScanLicenseText()
	if err != nil {
		return nil, err
	}
	retval := make([]string, 0, len(results[0].CycloneDXLicenses))
	for _, result := range results[0].CycloneDXLicenses {
		retval = append(retval, result.License.ID)
	}
	return retval, nil
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method, I only accept POST", http.StatusBadRequest)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/octet-stream" && contentType != "text/plain" {
		errorMessage := fmt.Sprintf("Invalid content type %s, I accept application/octet-stream or text/plain", contentType)
		http.Error(w, errorMessage, http.StatusBadRequest)
		return
	}

	content, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read file content", http.StatusInternalServerError)
		return
	}

	licenses, err := scanLicense(content)
	if err != nil {
		http.Error(w, "Failed to scan license", http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.MarshalIndent(licenses, "", "  ")
	if err != nil {
		http.Error(w, "Failed to marshal JSON response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Starting the web server on port %s...\n", port)

	http.HandleFunc("/", handleUpload)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s\n", err)
		return
	}
}
