package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/cenkalti/backoff/v5"
	"net/http"
	"net/url"
)

const (
	productServiceURL = "http://route256.pavl.uk:8080"
	token             = "testtoken"
)

type getProductRequest struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type GetProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

func getProduct(skuID int64) (g GetProductResponse, err error) {
	reqBody := getProductRequest{
		Token: token,
		SKU:   skuID,
	}
	reqBytes, err := json.Marshal(reqBody)
	if err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}
	url, err := url.JoinPath(productServiceURL, "get_product")
	if err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(reqBytes))
	if err != nil {
		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == 420 {
			return GetProductResponse{}, err
		}
		return GetProductResponse{}, backoff.Permanent(err)
	}
	defer resp.Body.Close()

	var productResp GetProductResponse

	if err := json.NewDecoder(resp.Body).Decode(&productResp); err != nil {
		return GetProductResponse{}, backoff.Permanent(err)
	}

	return GetProductResponse{Name: productResp.Name, Price: productResp.Price}, nil
}

func GetProduct(skuID int64) (g GetProductResponse, err error) {
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
