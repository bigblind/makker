package handler_helpers

import (
	"net/http"
	"encoding/json"
)

type errorResponse struct {
	Error error `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithJSONError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, errorResponse{err})
}