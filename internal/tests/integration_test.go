package tests

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"learn1/internal/client"
	h "learn1/internal/http"
	"learn1/internal/repo"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

type NotRealProductClient struct{}

func (f *NotRealProductClient) GetProduct(skuID int64) (client.GetProductResponse, error) {
	return client.GetProductResponse{
		Name:  "Сильфида",
		Price: 4510,
	}, nil
}

const productServiceURL = "http://localhost:3000"

func TestIntegrationGetCart(t *testing.T) {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})

	handler := h.NewHandler(r, &NotRealProductClient{})
	mux := handler.GetMux()
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

func TestIntegrationAddItem(t *testing.T) {
	r := repo.NewRepo()

	handler := h.NewHandler(r, &NotRealProductClient{})
	mux := handler.GetMux()
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

	addItem(t, 1, 5415913, 2)

	cart := getItems(t, 1)

	require.Len(t, cart.Items, 1)
	require.Equal(t, uint16(2), cart.Items[0].Count)
	require.Equal(t, uint32(9020), cart.TotalPrice)
}

func TestIntegrationDeleteItem(t *testing.T) {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{{SkuID: 5415913, Count: 1}})

	handler := h.NewHandler(r, &NotRealProductClient{})
	mux := handler.GetMux()
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

	deleteItem(t, 1, 5415913)

	cart := getItems(t, 1)

	require.Empty(t, cart.Items)
	require.Equal(t, uint32(0), cart.TotalPrice)
}

func TestIntegrationClearCart(t *testing.T) {
	r := repo.NewRepo()
	r.AddItem(1, []*repo.Item{
		{SkuID: 5415913, Count: 1},
		{SkuID: 123, Count: 2},
	})

	handler := h.NewHandler(r, &NotRealProductClient{})
	mux := handler.GetMux()
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

	clearCart(t, 1)

	cart := getItems(t, 1)

	require.Empty(t, cart.Items)
	require.Equal(t, uint32(0), cart.TotalPrice)
}

// Вспомогательные функции для тестов

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

func addItem(t *testing.T, userID, skuID, count int) {
	u := fmt.Sprintf("/user/%d/cart/%d", userID, skuID)
	url, _ := url.JoinPath(productServiceURL, u)

	body := fmt.Sprintf(`{"count": %d}`, count)
	resp, err := http.Post(url, "application/json", strings.NewReader(body))
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func deleteItem(t *testing.T, userID, skuID int) {
	u := fmt.Sprintf("/user/%d/cart/%d", userID, skuID)
	url, _ := url.JoinPath(productServiceURL, u)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}

func clearCart(t *testing.T, userID int) {
	u := fmt.Sprintf("/user/%d/cart", userID)
	url, _ := url.JoinPath(productServiceURL, u)

	req, _ := http.NewRequest(http.MethodDelete, url, nil)
	resp, err := http.DefaultClient.Do(req)

	require.NoError(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
}
