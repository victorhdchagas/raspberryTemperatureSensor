package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Server  ServerConfig  `json:"server"`
	DB      DBConfig      `json:"db"`
	Sensor  SensorConfig  `json:"sensor"`
	Weather WeatherConfig `json:"weather"`
}

type WeatherConfig struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
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
	gpioPin := "17"
	if envPin := os.Getenv("SENSOR_GPIO_PIN"); envPin != "" {
		gpioPin = envPin
	}

	// Penha, Rio de Janeiro (padrão). Sobrescreve via env:
	// WEATHER_LAT / WEATHER_LON.
	lat := -22.9094
	if envLat := os.Getenv("WEATHER_LAT"); envLat != "" {
		if f, err := strconv.ParseFloat(envLat, 64); err == nil {
			lat = f
		}
	}
	lon := -43.2746
	if envLon := os.Getenv("WEATHER_LON"); envLon != "" {
		if f, err := strconv.ParseFloat(envLon, 64); err == nil {
			lon = f
		}
	}

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
			GPIO:     gpioPin,
			Interval: 50 * time.Minute,
		},
		Weather: WeatherConfig{
			Lat: lat,
			Lon: lon,
		},
	}
}
