package core

import (
	"context"
	"net"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	emitterAPI "github.com/mingkaic/ultrasound/emitter/api"
	emitterProto "github.com/mingkaic/ultrasound/emitter/proto"
	viewerAPI "github.com/mingkaic/ultrasound/viewer/api"
	viewerProto "github.com/mingkaic/ultrasound/viewer/proto"
)

func Run(ctx context.Context, network, address string) error {
	log.Infof("GRPC listening on '%s'", address)
	listen, err := net.Listen(network, address)
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
		return err
	}
	defer func() {
		if err := listen.Close(); err != nil {
			log.Errorf("Failed to close %s %s: %v", network, address, err)
		}
	}()

	grpcServer := grpc.NewServer()
	viewerProto.RegisterViewerServer(grpcServer, viewerAPI.NewViewerServer())
	emitterProto.RegisterGraphEmitterServer(grpcServer, emitterAPI.NewEmitterServer())

	go func() {
		defer grpcServer.GracefulStop()
		<-ctx.Done()
	}()
	log.Infof("Starting listening at %s", address)
	return grpcServer.Serve(listen)
}
