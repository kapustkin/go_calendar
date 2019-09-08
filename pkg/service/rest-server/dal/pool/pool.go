package pool

import (
	"context"
	"fmt"
	"log"
	"time"

	grpcpool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
)

var (
	pool *grpcpool.Pool
)

// Init pool
func Init(addr string) {
	initPool, err := createGrpcConnectionPool(addr)
	if err != nil {
		log.Fatal(err)
	}

	pool = initPool
}

// GetConnectionFromPool return grpc connection
func GetConnectionFromPool(ctx *context.Context) (*grpcpool.ClientConn, error) {
	conn, err := pool.Get(*ctx)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func createGrpcConnectionPool(address string) (*grpcpool.Pool, error) {
	var factory grpcpool.Factory
	factory = func() (*grpc.ClientConn, error) {
		conn, err := grpc.Dial(address, grpc.WithInsecure())
		if err != nil {
			return nil, fmt.Errorf("Failed to start gRPC connection: %v", err)
		}
		return conn, err
	}
	pool, err := grpcpool.New(factory, 5, 5, 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("Failed to create gRPC pool: %v", err)
	}
	return pool, nil
}
