package http

import (
	"net/http"
)

type Server struct{}

func (s Server) Start() error {
	return http.ListenAndServe(":3000", NewRouter())
}
