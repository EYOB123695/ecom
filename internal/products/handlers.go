package products

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}

}

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products := []string{"hello", "world"}
	
	json.NewEncoder(w).Encode(products)
}