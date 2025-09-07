package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	jsonBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Println("Failed to marshal JSON:", err)
		return
	}

	w.Write(jsonBytes)

}

func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}
