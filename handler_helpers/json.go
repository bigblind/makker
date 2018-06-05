package handler_helpers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func RespondWithJSONError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, errorResponse{err.Error()})
}
