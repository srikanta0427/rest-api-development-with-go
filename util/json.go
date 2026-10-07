package util

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error"`
}

func JSONWriter(w http.ResponseWriter, status int, data Response) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status) // set the http code

	// into json
	return json.NewEncoder(w).Encode(data)
}

func Success(w http.ResponseWriter, status int, message string, data interface{}) {
	_ = JSONWriter(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(
	w http.ResponseWriter,
	status int,
	message string,
	err interface{},
) {
	_ = JSONWriter(w, status, Response{
		Success: false,
		Message: message,
		Error:   err,
	})
}
