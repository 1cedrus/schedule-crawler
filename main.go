package main

import (
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", scheduleHandler)
	mux.HandleFunc("/exam", examHandler)
	mux.HandleFunc("/{t}", scheduleHandler)

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("[ERROR]: Failed to listen on port 80: %v", err)
	}

	http.Serve(listener, mux)
}
