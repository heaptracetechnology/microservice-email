package http

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/oms-services/email/smtp"
)

type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler http.Handler
}

type Routes []Route

var routes = Routes{
	Route{
		"SendEmail",
		"POST",
		"/send",
		SendHandler{
			Emailer: smtp.Client{
				Address:  os.Getenv("SMTP_HOST") + ":" + os.Getenv("SMTP_PORT"),
				Password: os.Getenv("PASSWORD"),
			},
			Password: os.Getenv("PASSWORD"),
			SMTPHost: os.Getenv("SMTP_HOST"),
			SMTPPort: os.Getenv("SMTP_PORT"),
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

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		log.Println(route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.Handler)
	}

	return router
}
