package main

import (
	c "learn1/internal/client"
	"learn1/internal/db"
	h "learn1/internal/http"
	"learn1/internal/repo"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	sqldb, err := db.InitDB()
	if err != nil {
		panic(err)
	}
	defer sqldb.Close()
	/*
		r := repo.NewRepo()
		productClient := c.NewClient(c.ProductServiceURL, c.Token)
		r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})
		handler := h.NewHandler(r, productClient)
		mux := handler.GetMux()
		loggingMux := h.LoggingMiddleware(mux)
		http.ListenAndServe(":3000", loggingMux)
	*/
	r := repo.NewPostgresRepo(sqldb)
	r.GetItems(1)
	productClient := c.NewClient(c.ProductServiceURL, c.Token)
	handler := h.NewHandler(r, productClient)
	mux := handler.GetMux()
	loggingMux := h.LoggingMiddleware(mux)
	http.ListenAndServe(":3000", loggingMux)
}
