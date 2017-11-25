package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

// This method is the entry point for app engine.
func init() {
	router := mux.NewRouter()
	router.HandleFunc("/", readIndex)
	router.HandleFunc("/api/version", readVersion)
	http.Handle("/", router)
}

// API handlers

func readIndex(w http.ResponseWriter, r *http.Request) {

	// Assemble template path
	dir, err := os.Getwd()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	templatePath := filepath.Join(dir, "/template/index.html")
	t, err := template.New("index.html").ParseFiles(templatePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := map[string]string{"name": "world"}
	t.Execute(w, data)
}

func readVersion(w http.ResponseWriter, r *http.Request) {
	authors := []string{"Sandra Grujovic", "Paul Schmiedermayer", "Johannes Rohwer"}
	version := NewVersion("v0.1", authors)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(version)
}
