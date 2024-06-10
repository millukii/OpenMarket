package gateway

import (
	"context"
	"log"

	pb "github.com/millukii/commons/api"
	"github.com/millukii/commons/discovery"
	"github.com/millukii/commons/discovery/grpc"
)

type gateway struct {
	registry discovery.Registry
}

func NewGRPCGateway(registry discovery.Registry) *gateway{
	return &gateway{registry: registry}
}
func (g *gateway) CreateOrder(ctx context.Context,p *pb.CreateOrderRequest) (*pb.Order, error){
con, err := grpc.ServiceConnection(ctx, "orders", g.registry)
if err !=nil{
	log.Fatalf("Failed to dial server: %v", err)
	return nil, err
}
	c := pb.NewOrderServiceClient(con)

	return  c.CreateOrder(ctx, &pb.CreateOrderRequest{
			CustomerOd: p.CustomerOd ,
			Items: p.Items,
	})
}