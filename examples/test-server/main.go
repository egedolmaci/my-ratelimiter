package main

import (
	"net/http"
	"time"

	"github.com/egedolmaci/my-ratelimiter/internal/ratelimiter"
	"github.com/egedolmaci/my-ratelimiter/pkg/middleware"
)

type Rule int

type Server struct {
	mux        *http.ServeMux
	middleware middleware.Middleware
}

func NewServer() *Server {
	config := ratelimiter.Config{
		Strategy:   "fixed_window",
		Limit:      10,
		WindowSize: time.Minute,
	}

	rl := ratelimiter.NewRatelimiterWithConfig(config)
	defer rl.Stop()
	s := &Server{
		mux:        http.NewServeMux(),
		middleware: middleware.Middleware{Ratelimiter: rl},
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)

}

func (s *Server) routes() {
	s.mux.HandleFunc("/ratelimited", s.middleware.RateLimitMiddleware(UnlimitedHandler))
	s.mux.HandleFunc("/limited", LimitedHandler)
	s.mux.HandleFunc("/unlimited", UnlimitedHandler)
}

func main() {
	server := NewServer()
	http.ListenAndServe(":8080", server)
}
