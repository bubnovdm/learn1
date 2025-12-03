package main

import (
	h "learn1/internal/http"
	"learn1/internal/repo"
	"net/http"
)

func main() {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{{SkuID: 1, Count: 5}})
	mux := h.GetMux(r)
	http.ListenAndServe(":3000", mux)
}
