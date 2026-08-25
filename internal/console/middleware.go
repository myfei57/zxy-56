package console

import (
	"log"
	"net/http"
	"runtime/debug"
	"time"
)

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	})
}

func recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf("panic: %v\n%s", recovered, debug.Stack())
				writeError(w, http.StatusInternalServerError, errPanic)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
