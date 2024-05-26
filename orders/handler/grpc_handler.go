package handler

import (
	"context"
	"log"

	pb "github.com/millukii/commons/api"
	"github.com/millukii/openmarket-orders/service"
	"google.golang.org/grpc"

)
type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service service.Service
}

func NewGRPCHandler(grpcServer *grpc.Server, service service.Service) {
	handler := &grpcHandler{
		 service: service,
	}
	pb.RegisterOrderServiceServer(
		grpcServer,
		handler,
	)
}

func (h *grpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	log.Printf("New order... received %v",req )

	order := &pb.Order{
		ID:"12",
	}
	return order, nil
}