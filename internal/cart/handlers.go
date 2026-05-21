package cart

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
	"github.com/EYOB123695/ecom/internal/auth"
	appjson "github.com/EYOB123695/ecom/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

// Dummy usage to keep repo import for Swagger documentation
var _ = repo.CartItem{}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// GetCart godoc
// @Summary      Get cart
// @Tags         cart
// @Produce      json
// @Security     BearerAuth
// @Success      200 {object} CartResponse
// @Failure      401 {string} string
// @Failure      500 {string} string
// @Router       /cart [get]
func (h *handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cartResp, err := h.service.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, cartResp)
}

// AddCartItem godoc
// @Summary      Add item to cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        body body AddCartItemRequest true "Add cart item payload"
// @Success      201 {object} repo.CartItem
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /cart [post]
func (h *handler) AddCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.ProductID <= 0 {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}
	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	cartItem, err := h.service.AddCartItem(r.Context(), userID, req.ProductID, req.Quantity)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInsufficientStock) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusCreated, cartItem)
}

// UpdateCartItem godoc
// @Summary      Update cart item quantity
// @Tags         cart
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        product_id path int true "Product ID"
// @Param        body body UpdateCartItemRequest true "Update quantity payload"
// @Success      200 {object} repo.CartItem
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /cart/{product_id} [put]
func (h *handler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	prodIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.ParseInt(prodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than zero", http.StatusBadRequest)
		return
	}

	cartItem, err := h.service.UpdateCartItem(r.Context(), userID, productID, req.Quantity)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) || errors.Is(err, ErrCartItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if errors.Is(err, ErrInsufficientStock) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, cartItem)
}

// DeleteCartItem godoc
// @Summary      Remove item from cart
// @Tags         cart
// @Security     BearerAuth
// @Param        product_id path int true "Product ID"
// @Success      204 "No Content"
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /cart/{product_id} [delete]
func (h *handler) DeleteCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	prodIDStr := chi.URLParam(r, "product_id")
	productID, err := strconv.ParseInt(prodIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	err = h.service.DeleteCartItem(r.Context(), userID, productID)
	if err != nil {
		if errors.Is(err, ErrCartItemNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ClearCart godoc
// @Summary      Clear cart
// @Tags         cart
// @Security     BearerAuth
// @Success      204 "No Content"
// @Failure      401 {string} string
// @Failure      500 {string} string
// @Router       /cart [delete]
func (h *handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	err := h.service.ClearCart(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
