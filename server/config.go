package main

import (
	"flag"

	"github.com/mingkaic/ultrasound/data"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Network     string
	GRPCPort    int
	GatewayPort int
}

var (
	cfg      Config
	dbParams data.DBParams
)

func init() {
	flag.StringVar(&cfg.Network, "network", "tcp",
		`one of "tcp" or "unix". Must be consistent to -port`)
	flag.IntVar(&cfg.GatewayPort, "gw.port", 8080, "gateway port to serve on (default: 8080)")
	flag.IntVar(&cfg.GRPCPort, "grpc.port", 50051, "gRPC port to serve on (default: 50051)")
	dbParams.DeclFlags()
	flag.Parse()

	log.Infof("Server configuration: %+v", cfg)
	log.Infof("DB parameters: %+v", dbParams)
}
