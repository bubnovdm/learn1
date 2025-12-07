package http

import "learn1/internal/repo"

// Структура унифицированного овтета
type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

// Структура для POST, т.к. по условию coutn - body в запросе.
type addItemRequest struct {
	Count uint64 `json:"count"`
}

// Пока для ответа используется эта структура
type ItemsResponse struct {
	Items      []*repo.Item `json:"items"`
	TotalPrice uint32       `json:"total_price"`
}

// Наверное как-то так должна выглядель итоговая структура ответа (без суммы)
type CartItemResponse struct {
	SkuID int64  `json:"sku_id"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

// Итог с суммой
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalPrice uint32             `json:"total_price"`
}
