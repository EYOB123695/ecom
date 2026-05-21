package cart

// AddCartItemRequest defines the payload for adding an item to the cart
type AddCartItemRequest struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

// UpdateCartItemRequest defines the payload for updating item quantity
type UpdateCartItemRequest struct {
	Quantity int32 `json:"quantity"`
}

// CartItemResponse represents a cart item returned in the response
type CartItemResponse struct {
	ProductID    int64  `json:"product_id"`
	ProductName  string `json:"product_name"`
	ProductPrice int32  `json:"product_price"` // price in cents
	ProductStock int32  `json:"product_stock"` // remaining stock quantity
	Quantity     int32  `json:"quantity"`      // quantity in user's cart
}

// CartResponse represents the user's cart summary
type CartResponse struct {
	Items      []CartItemResponse `json:"items"`
	TotalCents int64              `json:"total_cents"`
}
