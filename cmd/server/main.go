package main

import (
	"fmt"
	"os"

	"github.com/oms-services/email/http"
)

func main() {
	getEnvOrExit("SMTP_HOST")
	getEnvOrExit("SMTP_PORT")
	getEnvOrExit("PASSWORD")

	server := http.Server{}
	if err := server.Start(); err != nil {
		panic(err)
	}
}

func getEnvOrExit(env string) string {
	value := os.Getenv(env)
	if value == "" {
		fmt.Fprintf(os.Stderr, "Environment variable '%s' must be set", env)
		os.Exit(1)
	}

	return value
}
