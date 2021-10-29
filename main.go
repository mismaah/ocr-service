package main

import (
	"encoding/json"
	"io"
	"log"
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
	router.HandleFunc("/imagetotext", ocrHandler).Methods("GET")
	log.Print("Server started on :" + webPort)
	err := http.ListenAndServe(":"+webPort, router)
	if err != nil {
		log.Print(err)
	}
}

func ocrHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
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
		log.Print(err)
		JSONError(w, "Unexpected error occured.", http.StatusInternalServerError)
		return
	}
	text, err := runOCR(resp.Body)
	if err != nil {
		JSONError(w, "Error reading from image.", http.StatusInternalServerError)
		return
	}
	result := map[string]interface{}{
		"text": text,
	}
	json.NewEncoder(w).Encode(&result)
}

func runOCR(image io.Reader) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()
	bytes, err := io.ReadAll(image)
	if err != nil {
		return "", err
	}
	client.SetImageFromBytes(bytes)
	text, err := client.Text()
	if err != nil {
		return "", err
	}
	return text, err
}

func JSONError(w http.ResponseWriter, errorMessage string, code int) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
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
		log.Fatalf("%s environment variable not set.\n", k)
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
