package products

import (

	"net/http"
	"github.com/EYOB123695/ecom/internal/json"
	
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
	err := h.service.ListProducts(r.Context())
	if err != nil {
		println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		

		return
	}
	products := []string{"hello", "world"}
	
	json.Write(w,http.StatusOK,products)

	
}