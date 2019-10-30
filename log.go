package main

import (
	"log"
	"net/http"
	"time"
)

func Log(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s [%s] %s %s ", r.RemoteAddr, r.Method, r.RequestURI, r.UserAgent())
		f(w, r)
	}
}

func LogResponseTime() func() {
	start := time.Now()
	return func() { log.Printf("%s\n", time.Since(start)) }
}
