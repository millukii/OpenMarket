package gateway

import (
	"context"

	pb "github.com/millukii/commons/api"

)

type OrdersGateway interface {
	CreateOrder(context.Context, *pb.CreateOrderRequest) (*pb.Order,error)
}