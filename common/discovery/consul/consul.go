package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	consul "github.com/hashicorp/consul/api"

)

type Registry struct {
	client *consul.Client
}

func NewRegistry(addr, serviceName string) (*Registry, error){
	config := consul.DefaultConfig()

	config.Address = addr
	client, err := consul.NewClient(config)

	if err != nil {
		return nil, err
	}

	return &Registry{
		client: client,
	},nil
}

	func (r *Registry) Register(ctx context.Context, instanceId, serviceName,hostPort string) error{
		parts := strings.Split(hostPort, ":")

		if len(parts) != 2 {
			return errors.New("invalid host:port format. Ex: localhost:8080")
		}

		port, err := strconv.Atoi(parts[1])
		if err !=nil{
			return err
		}

		host := parts[0]
		

		return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
			ID: instanceId,
			Address: host,
			Port: port,
			Name: serviceName,
			Check: &consul.AgentServiceCheck{
				CheckID: instanceId,
				TLSSkipVerify: true,
				TTL: "5s",
				Timeout: "1s",
				DeregisterCriticalServiceAfter: "10s",
			},
		})
	}
	
	func (r *Registry)  Deregister(ctx context.Context, instanceId, serverName string) error		{
		return r.client.Agent().CheckDeregister(instanceId)
	
	}
	
	func (r *Registry)  Discover(ctx context.Context, serviceName string) ([]string, error){

		entries, _ , err := r.client.Health().Service(serviceName, "", true, nil)
		if err !=nil{
			return nil,err
		}
		var instances []string

		for _,entry := range entries{
			instances = append(instances, fmt.Sprintf("%s:%d", entry.Service.Address, entry.Service.Port))
		}
		return instances,nil
	}
	
	func (r *Registry)  HealthCheck(instanceId, serviceName string) error{
		return r.client.Agent().UpdateTTL(instanceId,"online", api.HealthPassing)
	
	}