package service

import (
	"context"

	pb "github.com/millukii/commons/api"
	errors "github.com/millukii/commons/errors"
	"github.com/millukii/openmarket-orders/types"

)

type Service struct {
	store types.OrdersStore
}

func NewService(store types.OrdersStore) *Service {
	return &Service{store}
}

func (s *Service) CreateOrder(context.Context) error {
	return nil
}

func (s *Service)	ValidateOrder(ctx context.Context, r pb.CreateOrderRequest) error{
		if len(r.Items) ==0{
			return errors.ErrNoItems
		}
		return nil
}