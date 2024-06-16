package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	common "github.com/millukii/commons"
	"github.com/millukii/commons/broker/rabbitmq"
	"github.com/millukii/commons/discovery"
	"github.com/millukii/commons/discovery/consul"
	"github.com/millukii/openmarket-orders/handler"
	"github.com/millukii/openmarket-orders/service"
	"github.com/millukii/openmarket-orders/store"
	"google.golang.org/grpc"
)

var (
	serviceName = "orders"
	consulAddr = common.EnvString("CONSUL_ADDR", "localhost:8500")
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
	ampqHost = common.EnvString("RABBITMQ_HOST", "localhost")
	ampqPort = common.EnvString("RABBITMQ_PORT", "5672")
	ampqUser = common.EnvString("RABBITMQ_USER", "guest")
	ampqPass = common.EnvString("RABBITMQ_PASS", "guest")
)
func main() {

	ch, close :=	rabbitmq.Connect(ampqUser, ampqPass,ampqHost, ampqPort)

	defer func() {
		close()
		ch.Close()
	}()
	
	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)

	if  err !=nil {
		log.Fatalf("failed to listen %v",err)
	}
	defer l.Close()

 	store := store.NewStore()

 	svc := service.NewService(store)

 	handler.NewGRPCHandler(grpcServer, *svc, ch)

 	svc.CreateOrder(context.Background())

 	log.Println("GRPC Server started at ", grpcAddr)
	registry, err := consul.NewRegistry(consulAddr, serviceName)
		if err !=nil{
		log.Fatalf("Failed to registry gateway: %v", err)
	}
	
	ctx := context.Background()
	instanceId :=discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceId, serviceName,grpcAddr); err!=nil{
		log.Fatalf("Failed to registry: %v", err)
	}

	go func(){
		for {
			if err := registry.HealthCheck(instanceId, serviceName); err!=nil{
				log.Println(err)
				log.Fatal("failed to healthcheck")
			}
			time.Sleep(time.Second*1)
		}
	}()

	defer registry.Deregister(ctx, instanceId, serviceName)

	if err := grpcServer.Serve(l); err !=nil {
		log.Fatal(err.Error())
	}
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)
	go func() {
	<-term
	if err := l.Close(); err != nil {
		log.Fatalf("Error closing listener: %v", err)
	}
	grpcServer.Stop()
	}()
}