package middleware

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func logRequest(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request path: %s", r.URL.Path)
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RequestURI)

		defer func(startedAt time.Time) {
			log.Println(r.RequestURI, time.Since(startedAt))
		}(time.Now())

		next.ServeHTTP(w, r)
	})
}

func home(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}
