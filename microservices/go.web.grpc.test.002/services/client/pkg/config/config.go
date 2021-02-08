package config

import (
	"flag"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PDS_Host   string  `envconfig:"PDS_HOST"  default:"localhost:50051"`
	HTTP_Port  string  `envconfig:"HTTP_PORT" default:":5000"`
	InFile  string
}

func CreateConfig() (*Config, error){
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	flag.StringVar(&cfg.InFile, "in", "./ports.json", "Ports Json File")
	flag.Parse()

	return &cfg, nil
}
