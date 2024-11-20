package utils

import (
	"log"

	"gopkg.in/ini.v1"
)

type Config struct {
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	LogLevel       string
	MaxConnections int
	Timeout        int
}

func LoadConfig(path string) (*Config, error) {
	cfg, err := ini.Load(path)
	if err != nil {
		log.Printf("Error failed to load config file: %e", err)
		return nil, err
	}

	config := &Config{
		DBHost:         cfg.Section("database").Key("host").String(),
		DBPort:         cfg.Section("database").Key("port").MustInt(5432),
		DBUser:         cfg.Section("database").Key("user").String(),
		DBPassword:     cfg.Section("database").Key("password").String(),
		DBName:         cfg.Section("database").Key("dbname").String(),
		DBSSLMode:      cfg.Section("database").Key("sslmode").String(),
		LogLevel:       cfg.Section("app").Key("log_level").String(),
		MaxConnections: cfg.Section("app").Key("max_connections").MustInt(10),
		Timeout:        cfg.Section("app").Key("timeout").MustInt(30),
	}

	return config, nil
}
