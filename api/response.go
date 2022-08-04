package api

import (
	"encoding/json"
	"net/http"
	
	log "github.com/sirupsen/logrus"
)

type errorResponse struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, status int, err error) {
	log.WithError(err).Error("get configs")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	resp := errorResponse{err.Error()}
	response(w, status, resp)
}

func response(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	body, _ := json.Marshal(data)
	w.WriteHeader(status)
	_, _ = w.Write(body)
}
