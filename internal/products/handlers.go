package products

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EYOB123695/ecom/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}

}

// ListProducts godoc
// @Summary      List all products
// @Tags         products
// @Produce      json
// @Success      200 {array} ProductResponse
// @Failure      500 {string} string
// @Router       /products [get]
func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	json.Write(w, http.StatusOK, products)

}
// GetProductByID godoc
// @Summary      Get product by ID
// @Tags         products
// @Produce      json
// @Param        id path int true "Product ID"
// @Success      200 {object} ProductResponse
// @Failure      400 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /products/{id} [get]
func (h *handler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	product, err := h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, product)
}