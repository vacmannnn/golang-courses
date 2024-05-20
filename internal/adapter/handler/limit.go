package handler

import (
	"golang.org/x/time/rate"
	"net/http"
)

func limit(next http.Handler, bucketSize int) http.Handler {
	var limiter = rate.NewLimiter(1, bucketSize)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
