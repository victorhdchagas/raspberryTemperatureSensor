package sensor

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/d2r2/go-dht"
	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type Reading struct {
	Temp     float64
	Humidity float64
	Error    error
}

type DHT11 struct {
	pin int
}

func NewDHT11(pinName string) (*DHT11, error) {
	pin, err := strconv.Atoi(pinName)
	if err != nil {
		return nil, fmt.Errorf("invalid pin number: %w", err)
	}

	return &DHT11{
		pin: pin,
	}, nil
}

func (d *DHT11) Read() (temp, humidity float64, err error) {
	temperature32, humidity32, _, err := dht.ReadDHTxxWithRetry(dht.DHT11, d.pin, false, 10)
	if err != nil {
		return 0, 0, fmt.Errorf("sensor read error: %w", err)
	}

	if temperature32 < -40 || temperature32 > 80 {
		return 0, 0, fmt.Errorf("invalid temperature reading: %.2f", temperature32)
	}

	if humidity32 < 0 || humidity32 > 100 {
		return 0, 0, fmt.Errorf("invalid humidity reading: %.2f", humidity32)
	}

	return float64(temperature32), float64(humidity32), nil
}

func (d *DHT11) Start(ctx context.Context, interval time.Duration, database *db.Database) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Sensor collector started: reading every %v from GPIO pin %d", interval, d.pin)

	for {
		select {
		case <-ctx.Done():
			log.Println("Sensor collector stopped")
			return
		case <-ticker.C:
			temp, humidity, err := d.Read()
			if err != nil {
				log.Printf("Error reading sensor: %v", err)
				continue
			}

			if err := database.InsertMetric(temp, humidity); err != nil {
				log.Printf("Error inserting metric: %v", err)
				continue
			}

			log.Printf("Metric saved: Temp=%.1f°C, Humidity=%.1f%%", temp, humidity)
		}
	}
}

func (d *DHT11) Close() error {
	return nil
}
