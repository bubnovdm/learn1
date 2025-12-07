package http

import (
	"log"
	"net/http"
	"time"
)

/*
 Логгер из дополнительного п.1 из Readme
https://habr.com/ru/companies/otus/articles/857070/
*/

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("🚀 Старт обработки %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
		log.Printf("🏁 Завершено за %v", time.Since(start))
	})
}
