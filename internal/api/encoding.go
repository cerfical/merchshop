package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func writeResponse(w http.ResponseWriter, status int, r any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(r)
}

func readRequest[T any](r io.Reader) (*T, error) {
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()

	var req T
	if err := dec.Decode(&req); err != nil {
		return nil, err
	}

	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("request body must contain a single JSON object")
	}

	return &req, nil
}
