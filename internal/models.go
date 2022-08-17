package internal

type Product struct {
	Id          uint   `json:"id"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Price       uint   `json:"price"`
}
