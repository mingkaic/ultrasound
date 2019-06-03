package main

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/mingkaic/ultrasound/data"
	"github.com/mingkaic/ultrasound/server/core"
	"github.com/mingkaic/ultrasound/server/gateway"
)

func runServers(ctx context.Context) <-chan error {
	grpcAddr := fmt.Sprintf(":%d", cfg.GRPCPort)
	gwAddr := fmt.Sprintf(":%d", cfg.GatewayPort)
	ch := make(chan error, 2)
	// grpc service
	go func() {
		if err := core.Run(ctx, cfg.Network, grpcAddr); err != nil {
			ch <- fmt.Errorf("failed to run grpc service: %v", err)
		}
	}()
	// http gateway
	go func() {
		if err := gateway.Run(ctx, gateway.Options{
			Addr: gwAddr,
			GRPCServer: gateway.Endpoint{
				Network: cfg.Network,
				Addr:    grpcAddr,
			},
		}); err != nil {
			ch <- fmt.Errorf("failed to run gateway service: %v", err)
		}
	}()
	return ch
}

func main() {
	data.Open(&dbParams)
	defer data.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errCh := runServers(ctx)

	select {
	case err := <-errCh:
		log.Fatalf(err.Error())
		os.Exit(1)
	}
}
