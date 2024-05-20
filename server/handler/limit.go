package handler

import (
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"sync"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

// maybe inline it ?
func getVisitor(ip string, bucketLimit int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(1, bucketLimit)
		visitors[ip] = limiter
	}

	return limiter
}

func limit(next http.Handler, bucketLimit int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		limiter := getVisitor(ip, bucketLimit)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
