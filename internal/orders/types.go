package orders

// OrderItemResponse represents an item inside an order in the API response
type OrderItemResponse struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int32   `json:"quantity"`
	Price       float64 `json:"price"` // price of item in main currency (e.g. dollars)
}

// OrderResponse represents an order returned by the API
type OrderResponse struct {
	ID        int32               `json:"id"`
	UserID    int64               `json:"user_id"`
	Total     float64             `json:"total"`
	Status    string              `json:"status"`
	CreatedAt string              `json:"created_at"`
	Items     []OrderItemResponse `json:"items,omitempty"`
}
