package sensor

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type SensorManager struct {
	reader   *DHT11
	database *db.Database
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.RWMutex
}

func NewSensorManager(database *db.Database, gpioPin string) (*SensorManager, error) {
	reader, err := NewDHT11(gpioPin)
	if err != nil {
		return nil, err
	}

	return &SensorManager{
		reader:   reader,
		database: database,
	}, nil
}

func (sm *SensorManager) Start(interval time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.cancel != nil {
		sm.cancel()
	}

	sm.ctx, sm.cancel = context.WithCancel(context.Background())
	sm.reader.StartAsync(sm.ctx, interval, sm.database)
	log.Printf("Sensor manager started: reading every %v", interval)
}

func (sm *SensorManager) Restart(interval time.Duration) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.cancel != nil {
		sm.cancel()
	}

	sm.ctx, sm.cancel = context.WithCancel(context.Background())
	sm.reader.StartAsync(sm.ctx, interval, sm.database)
	log.Printf("Sensor manager restarted: reading every %v", interval)
}

func (sm *SensorManager) Stop() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.cancel != nil {
		sm.cancel()
		sm.cancel = nil
	}
	log.Println("Sensor manager stopped")
}
