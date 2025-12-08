package http

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

// Структура унифицированного овтета
type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}

// Структура для POST, т.к. по условию coutn - body в запросе.
type addItemRequest struct {
	Count int `json:"count"`
}

func (a addItemRequest) Validate() error {
	return validation.ValidateStruct(&a, validation.Field(&a.Count, validation.Required, validation.Min(1)))
}

type ValidURL struct {
	UserID int64 `json:"user_id"`
	SkuID  int64 `json:"sku_id"`
}

func (v ValidURL) Validate() error {
	return validation.ValidateStruct(&v,
		validation.Field(&v.UserID, validation.Min(1)),
		validation.Field(&v.SkuID, validation.Min(1)),
	)
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
