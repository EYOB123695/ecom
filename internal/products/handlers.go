package products

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	appjson "github.com/EYOB123695/ecom/internal/json"
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

	appjson.Write(w, http.StatusOK, products)

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

	appjson.Write(w, http.StatusOK, product)
}

type CreateProductRequest struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Quantity     int32  `json:"quantity"`
}

// CreateProduct godoc
// @Summary      Create a new product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        body body CreateProductRequest true "Create product payload"
// @Success      201 {object} ProductResponse
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /products [post]
func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "product name is required", http.StatusBadRequest)
		return
	}
	if req.PriceInCents <= 0 {
		http.Error(w, "price must be positive", http.StatusBadRequest)
		return
	}
	if req.Quantity < 0 {
		http.Error(w, "quantity cannot be negative", http.StatusBadRequest)
		return
	}

	product, err := h.service.CreateProduct(r.Context(), req.Name, req.PriceInCents, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusCreated, product)
}

type UpdateProductRequest struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Quantity     int32  `json:"quantity"`
}

// UpdateProduct godoc
// @Summary      Update a product
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id path int true "Product ID"
// @Param        body body UpdateProductRequest true "Update product payload"
// @Success      200 {object} ProductResponse
// @Failure      400 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /products/{id} [put]
func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "product name is required", http.StatusBadRequest)
		return
	}
	if req.PriceInCents <= 0 {
		http.Error(w, "price must be positive", http.StatusBadRequest)
		return
	}
	if req.Quantity < 0 {
		http.Error(w, "quantity cannot be negative", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product, err := h.service.UpdateProduct(r.Context(), id, req.Name, req.PriceInCents, req.Quantity)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary      Delete a product
// @Tags         products
// @Param        id path int true "Product ID"
// @Success      204 "No Content"
// @Failure      400 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /products/{id} [delete]
func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	_, err = h.service.GetProductByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "product not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}