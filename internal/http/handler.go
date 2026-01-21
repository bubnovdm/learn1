package http

import (
	"encoding/json"
	"learn1/internal/client"
	"learn1/internal/repo"
	"log"
	"net/http"
	"strconv"
)

func NewResponse(data any, err string) *Response {
	return &Response{Data: data, Error: err}
}

type Handler struct {
	repo   RepoInterface
	client client.ProductClient
}

func NewHandler(r RepoInterface, c client.ProductClient) *Handler {
	return &Handler{
		repo:   r,
		client: c,
	}
}

type ServerOptions struct {
	data       any
	statusCode int
	err        error
}

type OptParams func(*ServerOptions)

func WithData(data any) OptParams {
	return func(s *ServerOptions) {
		s.data = data
	}
}

func WithStatusCode(statusCode int) OptParams {
	return func(s *ServerOptions) {
		s.statusCode = statusCode
	}
}

func WithError(err error) OptParams {
	return func(s *ServerOptions) {
		s.err = err
	}
}

func WriteJSONResponse(w http.ResponseWriter, options ...OptParams) {
	server := &ServerOptions{
		statusCode: http.StatusOK,
	}
	for _, opt := range options {
		opt(server)
	}
	errorString := ""
	if server.err != nil {
		errorString = server.err.Error()
	}
	r := NewResponse(server.data, errorString)

	w.Header().Set("Content-Type", "application/json")

	resp, e := json.Marshal(r)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}

	w.WriteHeader(server.statusCode)
	w.Write(resp)
}

type RepoInterface interface {
	GetItems(userID int) ([]*repo.Item, error)
	AddItem(userID int, items []*repo.Item) error
	RemoveItem(userID int, skuID int) error
	ClearCart(userID int) error
}

func (h *Handler) GetMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/{user_id}/cart", h.getCart())
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", h.postAddItem())
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", h.deleteItem())
	mux.HandleFunc("DELETE /user/{user_id}/cart", h.deleteCart())
	return mux
}

func (h *Handler) getCart() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		itemsFromRepo, errdb := h.repo.GetItems(int(id))
		if errdb != nil {
			WriteJSONResponse(w, WithError(errdb))
		}
		// TODO: сортировка по skuID

		var cartResponse CartResponse
		var totalPrice uint32

		for _, item := range itemsFromRepo {
			sku := int64(item.SkuID)
			productResult, err := h.client.GetProduct(sku)
			if err != nil {
				WriteJSONResponse(w, WithError(err), WithStatusCode(http.StatusNotFound))
				return
			}

			cartResponse.Items = append(cartResponse.Items, CartItemResponse{
				SkuID: sku,
				Name:  productResult.Name,
				Count: item.Count,
				Price: productResult.Price})

			// Можно реализовать через пакет для денег
			totalPrice += productResult.Price * uint32(item.Count)
			cartResponse.TotalPrice = totalPrice

		}

		log.Printf("response: %v", cartResponse)
		WriteJSONResponse(w, WithData(cartResponse))
	}
}

func (h *Handler) postAddItem() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		skuIdString := req.PathValue("sku_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		skuId, err := strconv.ParseInt(skuIdString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		validUrl := ValidURL{UserID: id, SkuID: skuId}
		err = validUrl.Validate()
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		var body addItemRequest
		err = json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		err = body.Validate()
		if err != nil {
			WriteJSONResponse(w, WithError(err), WithStatusCode(http.StatusBadRequest))
			return
		}
		_, err = h.client.GetProduct(skuId)
		if err != nil {
			WriteJSONResponse(w,
				WithError(err),
				WithStatusCode(http.StatusNotFound))
			return
		}
		h.repo.AddItem(int(id), []*repo.Item{
			{SkuID: int(skuId), Count: uint16(body.Count)},
		})
		WriteJSONResponse(w)
	}
}

func (h *Handler) deleteItem() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		skuIdString := req.PathValue("sku_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		skuId, err := strconv.ParseInt(skuIdString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		h.repo.RemoveItem(int(id), int(skuId))
		WriteJSONResponse(w, WithStatusCode(http.StatusNoContent))
	}
}

func (h *Handler) deleteCart() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, WithError(err))
			return
		}
		h.repo.ClearCart(int(id))
		WriteJSONResponse(w, WithStatusCode(http.StatusNoContent))
	}
}
