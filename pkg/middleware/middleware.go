package middleware

import (
	"fmt"
	"net"
	"net/http"
)

type Limiter interface {
	IsRequestAllowed(identifier string) (bool, int) 
}

type Middleware struct {
	Ratelimiter Limiter

}

func (m *Middleware) RateLimitMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		host, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			panic("Error while spliting the identifier address")
		}
		identifier := host

		if allowed, remainingLimit := m.Ratelimiter.IsRequestAllowed(identifier); allowed {
			w.Write([]byte(fmt.Sprintf("Remaining limit = %d\n", remainingLimit)))
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	})
}
