package http

import (
	"net/http"
)

type Server struct {
	Routes []Route

	Emailer Emailer
}

func (s Server) Start() error {
	var routes = []Route{
		Route{
			"SendEmail",
			"POST",
			"/send",
			SendHandler{
				Emailer: s.Emailer,
			},
		},
		Route{
			"ReceiveEmail",
			"POST",
			"/receive",
			ReceiveHandler{},
		},
		Route{
			"Healthcheck",
			"Get",
			"/healthcheck",
			HealthcheckHandler{},
		},
	}

	return http.ListenAndServe(":3000", NewRouter(routes))
}
