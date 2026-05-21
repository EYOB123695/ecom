package products

import (
	"context"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	GetProductByID(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, name string, priceInCents int32, quantity int32) (repo.Product, error)
	UpdateProduct(ctx context.Context, id int64, name string, priceInCents int32, quantity int32) (repo.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

type svc struct {
	repo repo.Querier
}

func NewService(queries repo.Querier) Service {
	return &svc{repo: queries}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {

	return s.repo.ListProducts(ctx)

}

func (s *svc) GetProductByID(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.GetProductByID(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, name string, priceInCents int32, quantity int32) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         name,
		PriceInCents: priceInCents,
		Quantity:     quantity,
	})
}

func (s *svc) UpdateProduct(ctx context.Context, id int64, name string, priceInCents int32, quantity int32) (repo.Product, error) {
	return s.repo.UpdateProduct(ctx, repo.UpdateProductParams{
		ID:           id,
		Name:         name,
		PriceInCents: priceInCents,
		Quantity:     quantity,
	})
}

func (s *svc) DeleteProduct(ctx context.Context, id int64) error {
	return s.repo.DeleteProduct(ctx, id)
}
