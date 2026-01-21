package main

import (
	c "learn1/internal/client"
	h "learn1/internal/http"
	"learn1/internal/repo"
	"net/http"
)

func main() {
	r := repo.NewRepo()
	productClient := c.NewClient(c.ProductServiceURL, c.Token)
	r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})
	handler := h.NewHandler(r, productClient)
	mux := handler.GetMux()
	loggingMux := h.LoggingMiddleware(mux)
	http.ListenAndServe(":3000", loggingMux)
}
