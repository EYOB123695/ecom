package orders

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/EYOB123695/ecom/internal/auth"
	appjson "github.com/EYOB123695/ecom/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

// Checkout godoc
// @Summary      Checkout cart and place order
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Success      201 {object} OrderResponse
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      500 {string} string
// @Router       /orders [post]
func (h *handler) Checkout(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := h.service.Checkout(r.Context(), userID)
	if err != nil {
		if errors.Is(err, ErrCartEmpty) || errors.Is(err, ErrInsufficientStock) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusCreated, order)
}

// ListOrders godoc
// @Summary      List all user orders
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Success      200 {array} OrderResponse
// @Failure      401 {string} string
// @Failure      500 {string} string
// @Router       /orders [get]
func (h *handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.service.ListOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, orders)
}

// GetOrder godoc
// @Summary      Get order details by ID
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Param        id path int true "Order ID"
// @Success      200 {object} OrderResponse
// @Failure      400 {string} string
// @Failure      401 {string} string
// @Failure      404 {string} string
// @Failure      500 {string} string
// @Router       /orders/{id} [get]
func (h *handler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := auth.UserIDFromContext(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	orderID, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	order, err := h.service.GetOrderByID(r.Context(), int32(orderID), userID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	appjson.Write(w, http.StatusOK, order)
}
