package main

import (
	"net/http"
)

const (
	readinessPath  = "/healthz"
	fileServerPath = "/app/"
	rootPath       = "."
)

func main() {

	serveMux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir(rootPath))
	fS := http.StripPrefix(fileServerPath, fileServer)
	serveMux.Handle(fileServerPath, fS)

	serveMux.HandleFunc(readinessPath, handlerReadiness)

	server := http.Server{}
	server.Addr = ":8080"
	server.Handler = serveMux

	server.ListenAndServe()
}
