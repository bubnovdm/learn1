package client

import "fmt"

/*
 Тут будем обращаться к Product Service из п.2 Readme
 И тут же ретраи к Product Service из дополнительного п.3 Readme
*/

const (
	productServiceURL = "http://route256.pavl.uk:8080"
	token             = "testtoken"
)

type getProductRequest struct {
	Token string `json:"token"`
	SKU   int64  `json:"sku"`
}

type getProductResponse struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}

// Заглушка под внешний сервис
func TempResp(skuID int64) (name string, price uint32, err error) {
	switch skuID {
	case 1:
		return "Гречка ядрица", 49, nil
	case 2:
		return "Nike Air Zoom", 7999, nil
	case 3:
		return "Roxy Music. Stranded", 1028, nil
	case 4:
		return "Кроссовки Nike JORDAN", 2202, nil
	default:
		return "", 0, fmt.Errorf("product with sku %d not found", skuID)
	}
}

func GetProduct(skuID int64) (name string, price uint32, err error) {
	return
}
