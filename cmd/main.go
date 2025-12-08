package main

import (
	h "learn1/internal/http"
	"learn1/internal/repo"
	"net/http"
)

func main() {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})
	mux := h.GetMux(r)
	loggingMux := h.LoggingMiddleware(mux)
	http.ListenAndServe(":3000", loggingMux)
}
