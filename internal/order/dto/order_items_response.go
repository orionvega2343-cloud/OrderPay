package dto

type OrderItemResponse struct {
	Id           int    `json:"id"`
	OrderId      int    `json:"order_id"`
	ProductName  string `json:"product_name"`
	Quantity     int    `json:"quantity"`
	PricePerUnit int    `json:"price_per_unit"`
}
