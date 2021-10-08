package api

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Title   string `json:"title"`
	Details string `json:"details"`
}

func EncodeResponse(rw http.ResponseWriter, code int, body interface{}) error {
	rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	rw.WriteHeader(code)

	return json.NewEncoder(rw).Encode(body)
}
