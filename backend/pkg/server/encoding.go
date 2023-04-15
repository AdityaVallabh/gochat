package server

import (
	"encoding/json"
	"net/http"
)

func (s *Server) respond(w http.ResponseWriter, r *http.Request, status int, data any) {
	w.WriteHeader(status)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.Write([]byte("there was some error generating response\n"))
		}
	}
}

func (s *Server) decode(w http.ResponseWriter, r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
