package main

import log "github.com/sirupsen/logrus"
import "github.com/caarlos0/env/v6"

type db string

const (
	dbMemory db = "memory"
	dbEtcd3  db = "etcd3"
)

type config struct {
	DB         db       `env:"DB"`
	Etcd3Nodes []string `env:"ETCD3HOSTS" envSeparator:","`
	Address    string   `env:"ADDRESS"`
}

func getConfig() (cfg config) {
	cfg = config{}

	if err := env.Parse(&cfg); err != nil {
		log.Fatalf("%+v\n", err)
	}

	if cfg.DB == "" {
		log.Warn("Defaulting to in memory database")
		cfg.DB = dbMemory
	}

	if cfg.Address == "" {
		log.Info("Defaulting to :8080 address")
		cfg.Address = ":8080"
	}

	return
}
