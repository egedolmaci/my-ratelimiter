package main

import (
	"net/http"
	"time"

	"github.com/egedolmaci/my-ratelimiter/pkg/middleware"
)

type Rule int

type Server struct {
	mux *http.ServeMux
}

func NewServer() *Server {
	s := &Server{
		mux: http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)

}

func (s *Server) routes() {
	s.mux.HandleFunc("/ratelimited", middleware.RateLimitMiddleware(UnlimitedHandler, 1, time.Minute))
	s.mux.HandleFunc("/limited", LimitedHandler)
	s.mux.HandleFunc("/unlimited", UnlimitedHandler)
}

func main() {
	server := NewServer()
	http.ListenAndServe(":8080", server)
}
