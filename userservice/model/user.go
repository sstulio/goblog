package model

//User represents an user
type User struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	Products []Product `json:"products"`
	ServedBy string    `json:"servedBy"`
}

//Product represents a product
type Product struct {
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	ServedBy    string  `json:"servedBy"`
}
