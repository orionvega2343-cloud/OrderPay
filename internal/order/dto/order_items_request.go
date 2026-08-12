package dto

type OrderItemRequest struct {
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"price_per_unit"`
}
