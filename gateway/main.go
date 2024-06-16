package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	common "github.com/millukii/commons"
	"github.com/millukii/commons/discovery"
	"github.com/millukii/commons/discovery/consul"
	"github.com/millukii/openmarket-gateway/gateway"
	"github.com/millukii/openmarket-gateway/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":8080")
	consulAddr = common.EnvString("CONSUL_ADDR", "localhost:8500")
	orderService = "localhost:2000"
	serviceName = "gateway"
)

func main() {

	conn, err:= grpc.NewClient(orderService, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err !=nil{
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	registry , err := consul.NewRegistry(consulAddr, serviceName)
	
	if err !=nil{
		log.Fatalf("Failed to registry gateway: %v", err)
	}
	ctx := context.Background()
	instanceId :=discovery.GenerateInstanceID(serviceName)

	if err := registry.Register(ctx, instanceId, serviceName,httpAddr); err!=nil{
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
	
	ordersGateway := gateway.NewGRPCGateway(registry)

	log.Println("Dialing order service: ", orderService)
	mux := http.NewServeMux()

	handler := handlers.NewHttpHandler(ordersGateway)

	handler.RegisterRoutes(mux)

	log.Printf("Starting http server at %s", httpAddr)
	if err := http.ListenAndServe(httpAddr,mux); err !=nil{
		log.Fatal("failed to start")
	}

	srv := &http.Server{
	Addr: httpAddr,
	Handler: mux,
	}
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-term
		if err := srv.Close(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error closing Server: %v", err)
		}
	}()
}

