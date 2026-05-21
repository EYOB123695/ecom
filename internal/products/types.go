package products

// ProductResponse is used for API/Swagger documentation.
type ProductResponse struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Quantity     int32  `json:"quantity"`
	CreatedAt    string `json:"created_at"`
}
