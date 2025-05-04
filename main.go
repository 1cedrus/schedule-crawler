package main

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"no.name/qldt"
)

func handler(w http.ResponseWriter, r *http.Request) {
    tokenResp, err := qldt.FetchToken(r)
    if err != nil {
        fmt.Fprintf(w, "Error occured when fetching token!\n%v\n", err)
        return
    }

    fmt.Fprintf(w, "Token: %s\n", tokenResp.AccessToken)
}

func main() {
    mux := http.NewServeMux();
    mux.HandleFunc("/", handler);

    listener, err := net.Listen("tcp", ":80")
    if err != nil {
        log.Fatalf("[ERROR]: Failed to listen on port 80: %v", err)
    }

    http.Serve(listener, mux);
}
