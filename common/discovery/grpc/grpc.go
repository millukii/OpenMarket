package grpc

import (
	"context"
	"math/rand"

	"github.com/millukii/commons/discovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

)

func ServiceConnection(ctx context.Context, serviceName string, registry discovery.Registry)(*grpc.ClientConn,error){
	addrs ,err := registry.Discover(ctx, serviceName)

	if err != nil{
		return nil, err
	}

	return grpc.Dial(
		addrs[rand.Intn(len(addrs))],
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}