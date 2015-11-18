package main

import (
	"os"
	"strconv"
)

type Config struct {
	ConnectionString string
	SecretKey        string
	Port             string
	SessionExpires   int
}

func EnvWithDefault(name string, defaultVal string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultVal
	}
	return value
}

func Load() *Config {
	config := &Config{}

	ConnectionPort := EnvWithDefault("RDB_PORT_28015_TCP_PORT", "28015")
	ConnectionAddr := EnvWithDefault("RDB_PORT_28015_TCP_ADDR", "localhost")
	config.ConnectionString = ConnectionAddr + ":" + ConnectionPort
	config.SecretKey = EnvWithDefault("MAGNET_SESSION_KEY", "Here be dragons")
	config.Port = EnvWithDefault("MAGNET_PORT", ":3000")
	SessionExpires, err := strconv.Atoi(EnvWithDefault("MAGNET_SESSION_EXPIRE", "1296000"))
	if err != nil {
		config.SessionExpires = 1296000
	} else {
		config.SessionExpires = SessionExpires
	}

	return config
}
