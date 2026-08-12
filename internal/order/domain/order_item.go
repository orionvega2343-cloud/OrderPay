package domain

type OrderItem struct {
	Id           int    `db:"id"`
	OrderId      int    `db:"order_id"`
	ProductName  string `db:"product_name"`
	Quantity     int    `db:"quantity"`
	PricePerUnit int    `db:"price_per_unit"`
}
