package orders

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrCartEmpty         = errors.New("cart is empty")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrOrderNotFound     = errors.New("order not found")
)

type Service interface {
	Checkout(ctx context.Context, userID int64) (OrderResponse, error)
	ListOrders(ctx context.Context, userID int64) ([]OrderResponse, error)
	GetOrderByID(ctx context.Context, orderID int32, userID int64) (OrderResponse, error)
}

type svc struct {
	db   *pgx.Conn
	repo *repo.Queries
}

func NewService(db *pgx.Conn, queries *repo.Queries) Service {
	return &svc{db: db, repo: queries}
}

func (s *svc) Checkout(ctx context.Context, userID int64) (OrderResponse, error) {
	// Start Database Transaction
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	txQueries := s.repo.WithTx(tx)

	// 1. Fetch Cart Items
	cartItems, err := txQueries.GetCartItemsByUserID(ctx, userID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to fetch cart: %w", err)
	}
	if len(cartItems) == 0 {
		return OrderResponse{}, ErrCartEmpty
	}

	// 2. Validate Stock and Deduct Quantities
	var totalCents int64
	var orderItemsToCreate []repo.CreateOrderItemParams

	for _, item := range cartItems {
		// UpdateProductQuantity does: UPDATE products SET quantity = quantity - $2 WHERE id = $1 AND quantity >= $2 RETURNING *;
		// This handles the race conditions and stock checking in a single atomic SQL statement!
		updatedProduct, err := txQueries.UpdateProductQuantity(ctx, repo.UpdateProductQuantityParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			// If it returns pgx.ErrNoRows, it means the product doesn't exist or doesn't have enough stock!
			if errors.Is(err, pgx.ErrNoRows) {
				return OrderResponse{}, fmt.Errorf("%w for product %s (ID %d)", ErrInsufficientStock, item.ProductName, item.ProductID)
			}
			return OrderResponse{}, fmt.Errorf("failed to update stock for product %d: %w", item.ProductID, err)
		}

		itemPriceDollars := float64(updatedProduct.PriceInCents) / 100.0
		itemTotalCents := int64(updatedProduct.PriceInCents) * int64(item.Quantity)
		totalCents += itemTotalCents

		orderItemsToCreate = append(orderItemsToCreate, repo.CreateOrderItemParams{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     floatToNumeric(itemPriceDollars),
		})
	}

	totalDollars := float64(totalCents) / 100.0

	// 3. Create Order
	order, err := txQueries.CreateOrder(ctx, repo.CreateOrderParams{
		UserID: userID,
		Total:  floatToNumeric(totalDollars),
		Status: "completed", // orders processed successfully on checkout
	})
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to create order: %w", err)
	}

	// 4. Create Order Items
	var itemsResp []OrderItemResponse
	for i, itemParam := range orderItemsToCreate {
		itemParam.OrderID = order.ID
		oi, err := txQueries.CreateOrderItem(ctx, itemParam)
		if err != nil {
			return OrderResponse{}, fmt.Errorf("failed to create order item: %w", err)
		}

		itemsResp = append(itemsResp, OrderItemResponse{
			ProductID:   oi.ProductID,
			ProductName: cartItems[i].ProductName,
			Quantity:    oi.Quantity,
			Price:       numericToFloat(oi.Price),
		})
	}

	// 5. Clear Cart
	err = txQueries.ClearCart(ctx, userID)
	if err != nil {
		return OrderResponse{}, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit Transaction
	if err := tx.Commit(ctx); err != nil {
		return OrderResponse{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	createdAtStr := ""
	if order.CreatedAt.Valid {
		createdAtStr = order.CreatedAt.Time.Format(time.RFC3339)
	}

	return OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Total:     numericToFloat(order.Total),
		Status:    order.Status,
		CreatedAt: createdAtStr,
		Items:     itemsResp,
	}, nil
}

func (s *svc) ListOrders(ctx context.Context, userID int64) ([]OrderResponse, error) {
	orders, err := s.repo.ListOrdersByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	var resp []OrderResponse
	for _, order := range orders {
		createdAtStr := ""
		if order.CreatedAt.Valid {
			createdAtStr = order.CreatedAt.Time.Format(time.RFC3339)
		}

		resp = append(resp, OrderResponse{
			ID:        order.ID,
			UserID:    order.UserID,
			Total:     numericToFloat(order.Total),
			Status:    order.Status,
			CreatedAt: createdAtStr,
		})
	}

	return resp, nil
}

func (s *svc) GetOrderByID(ctx context.Context, orderID int32, userID int64) (OrderResponse, error) {
	order, err := s.repo.GetOrderByID(ctx, repo.GetOrderByIDParams{
		ID:     orderID,
		UserID: userID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return OrderResponse{}, ErrOrderNotFound
		}
		return OrderResponse{}, err
	}

	dbItems, err := s.repo.GetOrderItemsByOrderID(ctx, orderID)
	if err != nil {
		return OrderResponse{}, err
	}

	var items []OrderItemResponse
	for _, dbItem := range dbItems {
		items = append(items, OrderItemResponse{
			ProductID:   dbItem.ProductID,
			ProductName: dbItem.ProductName,
			Quantity:    dbItem.Quantity,
			Price:       numericToFloat(dbItem.Price),
		})
	}

	createdAtStr := ""
	if order.CreatedAt.Valid {
		createdAtStr = order.CreatedAt.Time.Format(time.RFC3339)
	}

	return OrderResponse{
		ID:        order.ID,
		UserID:    order.UserID,
		Total:     numericToFloat(order.Total),
		Status:    order.Status,
		CreatedAt: createdAtStr,
		Items:     items,
	}, nil
}

// Helpers for pgtype.Numeric conversion
func floatToNumeric(val float64) pgtype.Numeric {
	var num pgtype.Numeric
	_ = num.Scan(strconv.FormatFloat(val, 'f', 2, 64))
	return num
}

func numericToFloat(num pgtype.Numeric) float64 {
	if !num.Valid {
		return 0.0
	}
	f, _ := num.Float64Value()
	return f.Float64
}
