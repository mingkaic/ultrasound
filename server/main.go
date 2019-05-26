package main

import (
	"fmt"
	"net"

	log "github.com/sirupsen/logrus"

	"github.com/mingkaic/ultrasound/data"
	emitterAPI "github.com/mingkaic/ultrasound/emitter/api"
	emitterProto "github.com/mingkaic/ultrasound/emitter/proto"
	viewerAPI "github.com/mingkaic/ultrasound/viewer/api"
	viewerProto "github.com/mingkaic/ultrasound/viewer/proto"

	"google.golang.org/grpc"
)

func main() {
	data.Open(&dbParams)
	defer data.Close()

	host := fmt.Sprintf(":%d", cfg.Port)
	log.Infof("listening on '%s'", host)
	listen, err := net.Listen("tcp", host)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	viewerProto.RegisterViewerServer(grpcServer, viewerAPI.NewViewerServer())
	emitterProto.RegisterGraphEmitterServer(grpcServer, emitterAPI.NewEmitterServer())

	grpcServer.Serve(listen)
}
