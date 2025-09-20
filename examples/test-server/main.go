package testserver

import (
	"net/http"
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
	s.mux.HandleFunc("/unlimited", UnlimitedHandler)
	s.mux.HandleFunc("/limited", LimitedHandler)
}

func main() {
	server := &Server{}
	http.ListenAndServe(":8080", server)
}
