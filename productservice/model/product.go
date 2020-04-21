package model

//Product represents a product
type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:price`
	Description string  `json:description`
	ServedBy    string  `json:"servedBy"`
}
