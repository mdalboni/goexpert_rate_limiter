package middlewares

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/mdalboni/goexpert-rate-limiter/pkg/limiter"
)

type contextKey string

const rateLimitAllowedKey contextKey = "rateLimitAllowed"

// defaultKeyMapper returns the IP address of the request as the key for the rate limiter.
//
// It is used when no keyMapper function is provided to the NewRateLimiterMiddleware function.
// It returns an empty string if it fails to parse the IP address.
func defaultKeyMapper(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Printf("error parsing IP from RemoteAddr %s: %s\n", r.RemoteAddr, err)
		return ""
	}
	return fmt.Sprintf("IP:%s", ip)
}

// NewRateLimiterMiddleware returns a middleware that limits the number of requests per key.
//
// It uses the provided limiter.Limiter to check if the key has reached the limit.
//
// The keyMapper function is used to extract the key from the request.
// If the keyMapper function is nil, the defaultKeyMapper function is used.
//
// The middleware adds a context value "rateLimitAllowed" to the request context to avoid
// checking the rate limit for the same request multiple times, which allow for multiple rate limiters
// to be used in the same middleware chain.
func NewRateLimiterMiddleware(l *limiter.Limiter, keyMapper func(*http.Request) string) func(http.Handler) http.Handler {
	if keyMapper == nil {
		keyMapper = defaultKeyMapper
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := keyMapper(r)
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}
			if r.Context().Value(rateLimitAllowedKey) != true && l.ShouldLimit(key) {
				log.Printf("%s | LIMITED | %s", r.Context().Value(uidKey), key)
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte(`{"message": "you have reached the maximum number of requests or actions allowed within a certain time frame"}`))
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), rateLimitAllowedKey, true))
			next.ServeHTTP(w, r)
		})
	}
}
