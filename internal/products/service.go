package products

import (
	"context"
	repositories "github/budiharyonoo/ecom-tiago/internal/adapters/mysql/sqlc"
)

type Service interface {
	List(ctx context.Context) ([]repositories.Product, error)
	GetById(ctx context.Context, id uint64) (repositories.Product, error)
}

type svc struct {
	repo repositories.Querier
}

func NewService(r repositories.Querier) Service {
	return &svc{
		repo: r,
	}
}

func (s *svc) List(ctx context.Context) ([]repositories.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) GetById(ctx context.Context, id uint64) (repositories.Product, error) {
	return s.repo.GetProduct(ctx, id)
}
