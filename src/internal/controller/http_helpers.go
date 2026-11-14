package controller

import (
	"Notifications/src/internal/service"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type errorResponseDTO struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponseDTO{Error: message})
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrEmailAlreadyExists):
		writeError(w, http.StatusConflict, err.Error())
	case errors.Is(err, service.ErrInvalidCredentials):
		writeError(w, http.StatusUnauthorized, err.Error())
	case errors.Is(err, service.ErrNotFound):
		writeError(w, http.StatusNotFound, err.Error())
	case errors.Is(err, service.ErrUnsupportedChannel):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		log.Printf("unhandled service error: %v", err)
		writeError(w, http.StatusInternalServerError, "internal server error")
	}
}
