package client

/*
 Тут будем обращаться к Product Service из п.2 Readme
 И тут же ретраи к Product Service из дополнительного п.3 Readme
*/

type Product struct {
	Name  string `json:"name"`
	Price uint32 `json:"price"`
}
