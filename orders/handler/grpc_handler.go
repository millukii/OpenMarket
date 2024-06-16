package handler

import (
	"context"
	"encoding/json"
	"log"

	pb "github.com/millukii/commons/api"
	"github.com/millukii/commons/broker"
	"github.com/millukii/openmarket-orders/service"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)
type grpcHandler struct {
	pb.UnimplementedOrderServiceServer
	service service.Service
	ch *amqp.Channel
}

func NewGRPCHandler(grpcServer *grpc.Server, service service.Service,ch *amqp.Channel) {
	handler := &grpcHandler{
		 service: service,
		 ch: ch,
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
	marshalledOrder, err := json.Marshal(order)

	if err !=nil{
		return nil, err
	}
	q, err := h.ch.QueueDeclare(broker.OrderCreatedEvent,true,false,false,false,nil)

	if err !=nil{
		return nil, err
	}
	h.ch.PublishWithContext(ctx,"", q.Name, false,false,amqp.Publishing{
		ContentType: "application/json",
		Body: marshalledOrder,
		DeliveryMode: amqp.Persistent,
	})
	return order, nil
}