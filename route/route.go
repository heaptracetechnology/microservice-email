package route

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/oms-services/email/messaging"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"SendEmail",
		"POST",
		"/send",
		messaging.Send,
	},
	Route{
		"ReceiveEmail",
		"POST",
		"/receive",
		messaging.Receiver,
	},
	Route{
		"Healthcheck",
		"Get",
		"/healthcheck",
		healthcheck,
	},
}

func healthcheck(responseWriter http.ResponseWriter, request *http.Request) {
	responseWriter.WriteHeader(http.StatusOK)
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler
		log.Println(route.Name)
		handler = route.HandlerFunc

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}
	return router
}
