package products

import (
	"context"

	repo "github.com/EYOB123695/ecom/internal/adapters/postgresql/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	GetProductByID(ctx context.Context, id int64) (repo.Product, error)
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
