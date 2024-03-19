package middleware

import (
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Вызов следующего обработчика в цепочке
		next.ServeHTTP(w, r)

		// Логирование
		duration := time.Since(start)
		log.Printf("[%s] %s %s %s", r.Method, r.RemoteAddr, r.URL.Path, duration)
	})
}
