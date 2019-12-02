package http

import "net/http"

type HealthcheckHandler struct{}

func (h HealthcheckHandler) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
}
