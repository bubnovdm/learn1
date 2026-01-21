package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/cenkalti/backoff/v5"
)

const (
	ProductServiceURL = "http://route256.pavl.uk:8080"
	Token             = "testtoken"
)

type getProductRequest struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

type ProductClient interface {
	GetProduct(skuID int64) (GetProductResponse, error)
}

type Client struct {
	productServiceURL string
	token             string
}

func NewClient(productServiceURL string, token string) *Client {
	return &Client{
		productServiceURL: productServiceURL,
		token:             token,
	}
}

//TODO: вспомнить код
//TODO: перенести код GetProduct на структуру
//TODO: новый интерфейс для GetProduct
//TODO: handler только для интерфейса
//TODO: клиент передаётся через обертку
//TODO: в тесте реализация интерфейса, которая возвращает заглушку
//TODO: тест не ходит в реальный ProductService, а в заглушку
//TODO: остальные тесты

func getProduct(skuID int64) (g GetProductResponse, err error) {
	reqBody := getProductRequest{
		Token: Token,
		SKU:   skuID,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}
	url, err := url.JoinPath(ProductServiceURL, "get_product")
	if err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))

	if err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}
	// Чекаем на 429 и 420
	if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 420 {
		return GetProductResponse{}, err
	}
	// Чекаем на неок (типа 404), если получаем её, то не пытаемся ретраить а сразу даём перманентную ошибку
	if resp.StatusCode != http.StatusOK {
		//return GetProductResponse{}, backoff.Permanent(err)

		// Тут полупонятная херня, писал не сам)
		// Может как-то можно более аккуратно.
		var errResp map[string]any
		json.NewDecoder(resp.Body).Decode(&errResp)

		return GetProductResponse{}, backoff.Permanent(
			fmt.Errorf("Product Service returned %d: %v", resp.StatusCode, errResp),
		)
	}
	defer resp.Body.Close()

	var productResp GetProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}

	return GetProductResponse{Name: productResp.Name, Price: productResp.Price}, nil
}

func (Client) GetProduct(skuID int64) (g GetProductResponse, err error) {
	result, err := backoff.Retry(context.TODO(),
		func() (GetProductResponse, error) { return getProduct(skuID) },
		backoff.WithBackOff(backoff.NewExponentialBackOff()),
		backoff.WithMaxTries(3))
	if err != nil {
		fmt.Println("Error:", err)
		return GetProductResponse{}, err
	}
	return result, err
}
