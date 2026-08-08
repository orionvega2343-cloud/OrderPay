package dto

type OrderRequest struct {
	UserId string             `json:"user_id"`
	Items  []OrderItemRequest `json:"items"`
}
