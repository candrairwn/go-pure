package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func Encode[T any](w http.ResponseWriter, r *http.Request, status int, v T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}

// Wrapper function to add additional data
func EncodeWithWrapper[T any](w http.ResponseWriter, r *http.Request, status int, v T, additionalData map[string]interface{}) error {
	// Create a map to hold the original data and the additional data
	response := make(map[string]interface{})

	// Add the original data
	response["data"] = v

	// Add the additional data
	for key, value := range additionalData {
		response[key] = value
	}

	// Call the original Encode function with the wrapped response
	return Encode(w, r, status, response)
}
