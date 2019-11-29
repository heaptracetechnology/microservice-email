package main

import "github.com/oms-services/email/http"

func main() {
	server := http.Server{}
	if err := server.Start(); err != nil {
		panic(err)
	}
}
