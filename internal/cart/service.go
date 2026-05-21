package cart

import (
	"context"
	"errors"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
	ErrCartItemNotFound  = errors.New("cart item not found")
)

type Service interface {
	GetCart(ctx context.Context, userID int64) (CartResponse, error)
	AddCartItem(ctx context.Context, userID int64, productID int64, quantity int32) (repo.CartItem, error)
	UpdateCartItem(ctx context.Context, userID int64, productID int64, quantity int32) (repo.CartItem, error)
	DeleteCartItem(ctx context.Context, userID int64, productID int64) error
	ClearCart(ctx context.Context, userID int64) error
}

type svc struct {
	repo repo.Querier
}

func NewService(queries repo.Querier) Service {
	return &svc{repo: queries}
}

func (s *svc) GetCart(ctx context.Context, userID int64) (CartResponse, error) {
	dbItems, err := s.repo.GetCartItemsByUserID(ctx, userID)
	if err != nil {
		return CartResponse{}, err
	}

	var items []CartItemResponse
	var totalCents int64

	for _, item := range dbItems {
		itemTotal := int64(item.ProductPrice) * int64(item.Quantity)
		totalCents += itemTotal

		items = append(items, CartItemResponse{
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductPrice: item.ProductPrice,
			ProductStock: item.ProductStock,
			Quantity:     item.Quantity,
		})
	}

	return CartResponse{
		Items:      items,
		TotalCents: totalCents,
	}, nil
}

func (s *svc) AddCartItem(ctx context.Context, userID int64, productID int64, quantity int32) (repo.CartItem, error) {
	if quantity <= 0 {
		return repo.CartItem{}, errors.New("quantity must be positive")
	}

	// 1. Verify product exists
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrProductNotFound
		}
		return repo.CartItem{}, err
	}

	// 2. Fetch existing cart item to calculate new total quantity
	var currentQty int32
	dbItems, err := s.repo.GetCartItemsByUserID(ctx, userID)
	if err == nil {
		for _, item := range dbItems {
			if item.ProductID == productID {
				currentQty = item.Quantity
				break
			}
		}
	}

	newQty := currentQty + quantity
	if newQty > product.Quantity {
		return repo.CartItem{}, ErrInsufficientStock
	}

	// 3. Add or update cart item using upsert query
	return s.repo.AddCartItem(ctx, repo.AddCartItemParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	})
}

func (s *svc) UpdateCartItem(ctx context.Context, userID int64, productID int64, quantity int32) (repo.CartItem, error) {
	if quantity <= 0 {
		return repo.CartItem{}, errors.New("quantity must be positive")
	}

	// 1. Verify product exists
	product, err := s.repo.GetProductByID(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrProductNotFound
		}
		return repo.CartItem{}, err
	}

	if quantity > product.Quantity {
		return repo.CartItem{}, ErrInsufficientStock
	}

	// 2. Verify cart item exists
	var exists bool
	dbItems, err := s.repo.GetCartItemsByUserID(ctx, userID)
	if err == nil {
		for _, item := range dbItems {
			if item.ProductID == productID {
				exists = true
				break
			}
		}
	}

	if !exists {
		return repo.CartItem{}, ErrCartItemNotFound
	}

	return s.repo.UpdateCartItemQuantity(ctx, repo.UpdateCartItemQuantityParams{
		UserID:    userID,
		ProductID: productID,
		Quantity:  quantity,
	})
}

func (s *svc) DeleteCartItem(ctx context.Context, userID int64, productID int64) error {
	// Verify exists first
	var exists bool
	dbItems, err := s.repo.GetCartItemsByUserID(ctx, userID)
	if err == nil {
		for _, item := range dbItems {
			if item.ProductID == productID {
				exists = true
				break
			}
		}
	}

	if !exists {
		return ErrCartItemNotFound
	}

	return s.repo.DeleteCartItem(ctx, repo.DeleteCartItemParams{
		UserID:    userID,
		ProductID: productID,
	})
}

func (s *svc) ClearCart(ctx context.Context, userID int64) error {
	return s.repo.ClearCart(ctx, userID)
}
