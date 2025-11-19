package web

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	ErrorBody string `json:"error"`
}

func respondWithError(w http.ResponseWriter, code int, errorBody string) {
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Got here with issue %s", errorBody)

	errBody := errorResponse{
		ErrorBody: errorBody,
	}

	data, err := json.Marshal(errBody)
	if err != nil {
		log.Printf("Error marshalling error JSON: %s", err)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(payload)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to parse response JSON")
		return
	}

	w.WriteHeader(code)
	w.Write(res)
}
