package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"no.name/qldt"
)

const layout = "02/01/2006"

func handler(w http.ResponseWriter, r *http.Request) {
	t := r.PathValue("t")

	tokenResp, err := qldt.FetchToken(r)
	if err != nil {
		w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Error occured when fetching token!\n%v\n", err)
		return
	}

	scheduleResp, err := qldt.FetchDSTKB(tokenResp.AccessToken, tokenResp.Name)
    if err != nil {
        fmt.Fprintf(w, "Error occured when fetching schedule data!\n%v\n", err)
        return
    }

	now := time.Now()
	if len(t) > 0 {
		rt, err := strconv.Atoi(t)
		if err != nil {
			fmt.Fprintf(w, "Error occured when parse time!\n%v\n", err)
			return
		}

		now = now.Add(time.Duration(rt) * 7 * 24 * time.Hour)
	}

	var curTuanTKB qldt.TuanTKB
	for _, value := range scheduleResp.Data.DSTuanTKB {
		startTime, err := time.Parse(layout, value.NgayBatDau)
		if err != nil {
			fmt.Fprintf(w, "Error occured when parsing schedule data!\n%v\n", err)
            continue 
		}

		endTime, err := time.Parse(layout, value.NgayKetThuc)
		if err != nil {
			fmt.Fprintf(w, "Error occured when parsing schedule data!\n%v\n", err)
            continue
		}

		if now.After(startTime) && now.Before(endTime.Add(24*time.Hour-time.Nanosecond)) {
			curTuanTKB = value
			break
		}
	}

	fmt.Fprintf(w, "%v\n", curTuanTKB)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)
	mux.HandleFunc("/{t}", handler)

	listener, err := net.Listen("tcp", ":80")
	if err != nil {
		log.Fatalf("[ERROR]: Failed to listen on port 80: %v", err)
	}

	http.Serve(listener, mux)
}
