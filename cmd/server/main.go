package main

import (
	"fmt"
	"os"

	"github.com/oms-services/email/http"
	"github.com/oms-services/email/smtp"
)

func main() {
	smtpHost := getEnvOrExit("SMTP_HOST")
	smtpPort := getEnvOrExit("SMTP_PORT")
	password := getEnvOrExit("PASSWORD")

	smtpClient := smtp.Client{
		Address:  smtpHost + ":" + smtpPort,
		Password: password,
	}

	server := http.Server{
		Emailer: smtpClient,
	}

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
