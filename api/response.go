package api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, status int, message string) {
	w.Header().Add("Content-Type", "application/json")
	resp := ErrorResponse{message}
	response(w, status, resp)
}

func response(w http.ResponseWriter, status int, data any) {
	w.Header().Add("Content-Type", "application/json")
	// TODO: Enhance this
	w.Header().Add("Access-Control-Allow-Origin", "*")
	body, _ := json.Marshal(data)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
