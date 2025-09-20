package main

import "net/http"

func UnlimitedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Unlimited endpoint - no rate limiting"))

}

func LimitedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusTooManyRequests)
	w.Write([]byte("Limited endpoint - it is rate limited"))
}
