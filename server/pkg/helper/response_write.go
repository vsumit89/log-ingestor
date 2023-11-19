package helper

import (
	"encoding/json"
	"net/http"
)

type jsonResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func SendJSONResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	response := jsonResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON response", http.StatusInternalServerError)
	}
}
