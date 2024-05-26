package main

import (
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	common "github.com/millukii/commons"
	pb "github.com/millukii/commons/api"
	"github.com/millukii/openmarket-gateway/handlers"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

)

var (
	httpAddr = common.EnvString("HTTP_ADDR", ":8080")
	orderService = "localhost:2000"
)

func main() {

	conn, err:= grpc.Dial(orderService, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err !=nil{
		log.Fatalf("Failed to dial server: %v", err)
	}

	defer conn.Close()

	c := pb.NewOrderServiceClient(conn)

	log.Println("Dialing order service: ", orderService)
	mux := http.NewServeMux()

	handler := handlers.NewHttpHandler(c)

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

