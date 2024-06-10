package discovery

import (
	"context"
)

type Registry interface {
	Register(ctx context.Context, instanceId, serverName,hostPort string) error
	Deregister(ctx context.Context, instanceId, serverName string) error
	Discover(ctx context.Context, serviceName string) ([]string, error)
	HealthCheck(instanceId, serviceName string) error
}