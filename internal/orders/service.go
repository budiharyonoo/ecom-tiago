package orders

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	repositories "github/budiharyonoo/ecom-tiago/internal/adapters/mysql/sqlc"
)

var (
	ErrProductNotFound = errors.New("product not found")
	ErrProductNoStock  = errors.New("product has not enough stock")
)

type Service interface {
	Store(ctx context.Context, orderRequest CreateOrderRequest) (sql.Result, error)
}

type svc struct {
	repo *repositories.Queries
	db   *sql.DB
}

func NewService(r *repositories.Queries, db *sql.DB) Service {
	return &svc{
		repo: r,
		db:   db,
	}
}

func (s *svc) Store(ctx context.Context, orderRequest CreateOrderRequest) (sql.Result, error) {
	// begin DB transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	qtx := s.repo.WithTx(tx)

	// insert to order table
	result, err := qtx.CreateOrder(ctx, repositories.CreateOrderParams{
		CustomerID: orderRequest.CustomerID,
		TotalPrice: orderRequest.TotalPrice,
		Status:     "ordered",
	})
	if err != nil {
		return nil, err
	}

	orderId, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	// validate to products to check existing product_id
	for _, item := range orderRequest.Items {
		product, err := qtx.GetProduct(ctx, uint64(item.ProductID))
		if err != nil {
			return nil, fmt.Errorf("%w: ID %d", ErrProductNotFound, item.ProductID)
		}

		if item.Quantity > product.Quantity {
			return nil, fmt.Errorf("%w: ID %d", ErrProductNoStock, item.ProductID)
		}

		// insert for order_items on the loop because sqlc doesnt support bulk insert,
		// but no worries because using TX it will only commited to the disk once
		_, err = qtx.CreateOrderItems(ctx, repositories.CreateOrderItemsParams{
			OrderID:     orderId,
			ProductID:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    item.Quantity,
		})
		if err != nil {
			return nil, err
		}

		// update the stock qty
		_, err = qtx.UpdateProductQty(ctx, repositories.UpdateProductQtyParams{
			ID:       uint64(item.ProductID),
			Quantity: item.Quantity,
		})
		if err != nil {
			return nil, err
		}
	}

	tx.Commit()

	return result, nil
}
