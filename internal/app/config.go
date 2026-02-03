package app

import (
	"sync"
	"time"
)

type Config struct {
	SensorReadingInterval time.Duration
	mu                    sync.RWMutex
}

func NewConfig() *Config {
	return &Config{
		SensorReadingInterval: 50 * time.Minute,
	}
}

func (c *Config) GetSensorReadingInterval() time.Duration {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.SensorReadingInterval
}

func (c *Config) SetSensorReadingInterval(interval time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SensorReadingInterval = interval
}
