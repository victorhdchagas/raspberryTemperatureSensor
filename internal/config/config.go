package config

import (
	"time"
)

type Config struct {
	Server ServerConfig `json:"server"`
	DB     DBConfig     `json:"db"`
	Sensor SensorConfig `json:"sensor"`
}

type ServerConfig struct {
	Port         string        `json:"port"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

type DBConfig struct {
	Path string `json:"path"`
}

type SensorConfig struct {
	GPIO     string        `json:"gpio"`
	Interval time.Duration `json:"interval"`
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:         ":8080",
			ReadTimeout:  15 * time.Second,
			WriteTimeout: 15 * time.Second,
		},
		DB: DBConfig{
			Path: "./data/telemetry.db",
		},
		Sensor: SensorConfig{
			GPIO:     "4",
			Interval: 50 * time.Minute,
		},
	}
}
