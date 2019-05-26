package main

import (
	"flag"

	"github.com/mingkaic/ultrasound/data"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Port int
}

var (
	cfg      Config
	dbParams data.DBParams
)

func init() {
	flag.IntVar(&cfg.Port, "port", 8080, "gRPC port to serve on (default: 8080)")
	dbParams.DeclFlags()
	flag.Parse()

	log.Infof("Server configuration: %+v", cfg)
	log.Infof("DB parameters: %+v", dbParams)
}
