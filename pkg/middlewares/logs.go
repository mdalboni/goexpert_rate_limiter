package middlewares

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
)

const uidKey contextKey = "uid"

// LogRequest logs the request method, path and IP address of the request.
func LogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uid := uuid.New().String()
		r = r.WithContext(context.WithValue(r.Context(), uidKey, uid))

		log.Printf(
			"Start | %s | %s %s | IP %s",
			uid,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)
		next.ServeHTTP(w, r)
		log.Printf(
			"End | %s | %s %s | IP %s",
			uid,
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
		)
	})
}
