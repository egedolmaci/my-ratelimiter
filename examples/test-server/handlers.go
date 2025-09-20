package testserver

import "net/http"

func UnlimitedHandler(w http.ResponseWriter, r *http.Request) {

}

func LimitedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(429)
}
