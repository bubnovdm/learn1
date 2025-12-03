package http

import (
	"encoding/json"
	"learn1/internal/repo"
	"net/http"
	"strconv"
)

type Response struct {
	Data  any    `json:"data"`
	Error string `json:"error,omitempty"`
}
type addItemRequest struct {
	Count uint64 `json:"count"`
}

func NewResponse(data any, err string) *Response {
	return &Response{Data: data, Error: err}
}

func WriteJSONResponse(w http.ResponseWriter, data any, err error) {
	errorString := ""
	if err != nil {
		errorString = err.Error()
	}
	r := NewResponse(data, errorString)
	w.Header().Set("Content-Type", "application/json")
	ans, e := json.Marshal(r)
	if e != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(e.Error()))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(ans)
}

type RepoInterface interface {
	GetItems(userID int) []*repo.Item
	AddItem(userID int, items []*repo.Item)
	RemoveItem(userID int, skuID int)
	ClearCard(userID int)
}

func GetMux(r RepoInterface) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/{user_id}/cart", getCart(r))
	mux.HandleFunc("POST /user/{user_id}/cart/{sku_id}", postAddItem(r))
	mux.HandleFunc("DELETE /user/{user_id}/cart/{sku_id}", deleteItem(r))
	mux.HandleFunc("DELETE /user/{user_id}/cart", deleteCart(r))
	return mux
}

type ItemsResponse struct {
	Items []*repo.Item `json:"items"`
}

func getCart(r RepoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		data := r.GetItems(int(id))
		items := ItemsResponse{Items: data}
		WriteJSONResponse(w, items, nil)
	}
}

func postAddItem(r RepoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		skuIdString := req.PathValue("sku_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		skuId, err := strconv.ParseInt(skuIdString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		var body addItemRequest
		err = json.NewDecoder(req.Body).Decode(&body)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		r.AddItem(int(id), []*repo.Item{
			{SkuID: int(skuId), Count: int(body.Count)},
		})
		WriteJSONResponse(w, nil, nil)
	}
}

func deleteItem(r RepoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		skuIdString := req.PathValue("sku_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		skuId, err := strconv.ParseInt(skuIdString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		r.RemoveItem(int(id), int(skuId))
		WriteJSONResponse(w, nil, nil)
	}
}

func deleteCart(r RepoInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		idString := req.PathValue("user_id")
		id, err := strconv.ParseInt(idString, 10, 64)
		if err != nil {
			WriteJSONResponse(w, nil, err)
			return
		}
		r.ClearCard(int(id))
		w.WriteHeader(http.StatusNoContent)
	}
}
