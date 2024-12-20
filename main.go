package main

import (
	"net/http"
)

func main() {

	serveMux := http.ServeMux{}

	server := http.Server{}
	server.Addr = ":8080"
	server.Handler = &serveMux

	server.ListenAndServe()
}
