package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	common "github.com/millukii/commons"
	"github.com/millukii/openmarket-orders/handler"
	"github.com/millukii/openmarket-orders/service"
	"github.com/millukii/openmarket-orders/store"
	"google.golang.org/grpc"
)

var (
	grpcAddr = common.EnvString("GRPC_ADDR", "localhost:2000")
)
func main() {

	grpcServer := grpc.NewServer()

	l, err := net.Listen("tcp", grpcAddr)

	if  err !=nil {
		log.Fatalf("failed to listen %v",err)
	}
	defer l.Close()

 store := store.NewStore()

 svc := service.NewService(store)

 handler.NewGRPCHandler(grpcServer, *svc)

 svc.CreateOrder(context.Background())

 log.Println("GRPC Server started at ", grpcAddr)

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