package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PDS_Port  string  `envconfig:"PDS_PORT" default:":50051"`
	SQL_Host  string  `envconfig:"SQL_HOST" default:"localhost:15432"`
	SQL_User  string  `envconfig:"SQL_USER" default:"postgres"`
	SQL_Pass  string  `envconfig:"SQL_PASS" default:"mysecretpassword"`
}

func CreateConfig() (*Config, error){
	var env Config
	err := envconfig.Process("", &env)
	if err != nil {
		return nil, err
	}
	return &env, nil
}
