package service

import (
	"context"
	"log"

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
	if (len(r.Items) ==0){
		return errors.ErrNoItems
	}

	for _,i :=range r.Items{
		if i.ID == ""{
				return errors.ErrNoId
		}
		if i.Quantity <= 0{
				return errors.ErrInvalidQuantity
		}
	}
	mergedItems := mergeItemsQuantities(r.Items)

	log.Println("merged items ", mergedItems)
	//validate with stock service
	return nil
}

func mergeItemsQuantities(items []*pb.ItemsWithQuantity) []*pb.ItemsWithQuantity{
	merged := make([]*pb.ItemsWithQuantity,0)

	for _, item := range items{
		found := false

		for _,finalItem := range merged {
			if finalItem.ID == item.ID{
				finalItem.Quantity += item.Quantity
				found = true
				break
			}
			if !found{
				merged = append(merged, item)
			}
		}
	}

	return merged
}