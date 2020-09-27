package main

import log "github.com/sirupsen/logrus"
import "github.com/caarlos0/env/v6"

type db string

const (
	dbMemory db = "memory"
	dbEtcd3     = "etcd3"
	dbDisk      = "disk"
	dbRedis     = "redis"
	dbBadger    = "badger"
	dbNull      = "null"
)

type config struct {
	DB       db       `env:"DB"`
	DBHosts  []string `env:"DBHOST" envSeparator:","`
	DiskPath string   `env:"DISKPATH"`
	Address  string   `env:"ADDRESS"`
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

	if cfg.DB == dbDisk && cfg.DiskPath == "" {
		log.Warn("Defaulting to ./data path for disk storage")
		cfg.DiskPath = "./data"
	}

	if cfg.DB == dbNull {
		log.Warn("Are you sure?")
	}

	if len(cfg.DBHosts) == 0 && (cfg.DB == dbRedis || cfg.DB == dbEtcd3) {
		log.Fatalf("No database host(s) specified but was required")
	}

	if cfg.Address == "" {
		log.Info("Defaulting to :8080 address")
		cfg.Address = ":8080"
	}

	return
}
