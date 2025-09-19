package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/unlimited", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Access to this resource is limited."))
	})

	http.ListenAndServe(":8080", nil)
}