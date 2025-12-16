package tests

import (
	"encoding/json"
	"fmt"
	h "learn1/internal/http"
	"learn1/internal/repo"
	"log"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIntegrationCart(t *testing.T) {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})
	mux := h.GetMux(r)
	loggingMux := h.LoggingMiddleware(mux)
	srv := &http.Server{Addr: ":3000", Handler: loggingMux}
	defer srv.Close()
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
			return
		}
	}()
	cartResponse := getItems(t, 1)
	expected := &h.CartResponse{
		Items: []h.CartItemResponse{{
			SkuID: 5415913,
			Name:  "Сильфида",
			Count: 1,
			Price: 4510,
		}},
		TotalPrice: 4510,
	}
	require.Equal(t, expected, cartResponse)
}

const productServiceURL = "http://localhost:3000"

func getItems(t *testing.T, userID int) *h.CartResponse {
	u := fmt.Sprintf("/user/%d/cart", userID)
	url, err := url.JoinPath(productServiceURL, u)
	if err != nil {
		t.Fatalf("could not create url: %v", err)
		return nil
	}

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("could not get cart: %v", err)
	}

	// Чекаем на неок (типа 404), если получаем её, то не пытаемся ретраить а сразу даём перманентную ошибку
	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)

		t.Fatalf("could not get cart: %v, status code: %v", errResp, resp.StatusCode)
		return nil
	}
	defer resp.Body.Close()

	type serverResponse struct {
		Data  h.CartResponse `json:"data"`
		Error string         `json:"error,omitempty"`
	}
	var server serverResponse
	err = json.NewDecoder(resp.Body).Decode(&server)
	if err != nil {
		t.Fatalf("could not decode cart: %v", err)
		return nil
	}

	if server.Error != "" {
		t.Fatalf("cart: %v", server.Error)
	}
	return &server.Data
}
