package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

var logger = log.New(os.Stdout, "http: ", log.LstdFlags)

//Error stores authentication error information and handles response
type Error struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func writeResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
