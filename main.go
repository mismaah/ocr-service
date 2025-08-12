package main

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/otiai10/gosseract/v2"
)

var (
	webPort string
	authKey string
)

func main() {
	loadEnv()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/health", healthCheckHandler).Methods("GET")
	router.HandleFunc("/text/url", urlHandler).Methods("POST")
	router.HandleFunc("/text/upload", uploadHandler).Methods("POST")
	slog.Info("Server started", "PORT", webPort)
	err := http.ListenAndServe(":"+webPort, router)
	if err != nil {
		slog.Error("Starting server", "error", err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	validToken := isTokenValid(w, r)
	if !validToken {
		JSONError(w, "Invalid API key.", http.StatusUnauthorized)
		return
	}
	file, _, err := r.FormFile("file")
	if err != nil {
		JSONError(w, "Error uploading file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	text, err := runOCR(file)
	if err != nil {
		JSONError(w, "Error reading from image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := map[string]string{
		"text": text,
	}
	json.NewEncoder(w).Encode(&result)
}

func urlHandler(w http.ResponseWriter, r *http.Request) {
	validToken := isTokenValid(w, r)
	if !validToken {
		JSONError(w, "Invalid token.", http.StatusUnauthorized)
		return
	}
	params := r.URL.Query()
	url := params.Get("url")
	if url == "" {
		JSONError(w, "Image URL not given.", http.StatusBadRequest)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		JSONError(w, "Error retrieving image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	text, err := runOCR(resp.Body)
	if err != nil {
		JSONError(w, "Error reading from image: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	result := map[string]string{
		"text": text,
	}
	json.NewEncoder(w).Encode(&result)
}

func runOCR(image io.Reader) (string, error) {
	client := gosseract.NewClient()
	// client.SetWhitelist("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789.,;:!?-()[]\"'`~@#$%^&*+=_\\|/<> ")
	defer client.Close()
	bytes, err := io.ReadAll(image)
	if err != nil {
		return "", errors.New("error reading image data: " + err.Error())
	}
	err = client.SetImageFromBytes(bytes)
	if err != nil {
		return "", errors.New("error setting image from bytes: " + err.Error())
	}
	text, err := client.Text()
	if err != nil {
		return "", errors.New("error extracting text from image: " + err.Error())
	}
	return text, err
}

func JSONError(w http.ResponseWriter, errorMessage string, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]string{
		"error": errorMessage,
	}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}

func loadEnv() {
	godotenv.Load()
	webPort = mustGetenv("PORT")
	authKey = os.Getenv("AUTH_KEY")
}

func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		slog.Error("Getting environment variable", "variable", k)
	}
	return v
}

func isTokenValid(w http.ResponseWriter, r *http.Request) bool {
	if authKey == "" {
		return true
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == authKey {
		return true
	}
	return false
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"status": "ok"}
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		JSONError(w, "Error encoding response: "+err.Error(), http.StatusInternalServerError)
	}
}
