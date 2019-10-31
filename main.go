package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr:         "0.0.0.0:80",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	http.HandleFunc("/", Log(HomeHandler))
	http.HandleFunc("/help", Log(InfoHandler))
	http.HandleFunc("/request", Log(ClientRequestHandler))
	http.HandleFunc("/memory", Log(SystemHandler))
	http.HandleFunc("/return", Log(ReturnCodeHandler))
	http.HandleFunc("/headers", Log(ReturnHeadersHandler))
	http.HandleFunc("/health", Log(HealthCheckHandler))
	http.HandleFunc("/interfaces", Log(InterfacesHandler))
	log.Fatal(s.ListenAndServe())
}
