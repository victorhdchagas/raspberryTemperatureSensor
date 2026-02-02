package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
)

type Reading struct {
	Temp     float64 `json:"temperature"`
	Humidity float64 `json:"humidity"`
	Error    *string `json:"error"`
}

type DHT11 struct {
	pinName string
}

func NewDHT11(pinName string) (*DHT11, error) {
	// We no longer need periph.io here.
	return &DHT11{
		pinName: pinName,
	}, nil
}

func (d *DHT11) Read() (temp, humidity float64, err error) {
	// Execute the Python script
	cmd := exec.Command("python3", "scripts/read_sensor.py", d.pinName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return 0, 0, fmt.Errorf("failed to execute python script: %v (output: %s)", err, string(output))
	}

	var reading Reading
	if err := json.Unmarshal(output, &reading); err != nil {
		return 0, 0, fmt.Errorf("failed to parse python output: %v (output: %s)", err, string(output))
	}

	if reading.Error != nil && *reading.Error != "" {
		return 0, 0, fmt.Errorf("sensor error: %s", *reading.Error)
	}

	return reading.Temp, reading.Humidity, nil
}

func (d *DHT11) Start(ctx context.Context, interval time.Duration, database *db.Database) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Sensor collector started (via Python wrapper): reading every %v from GPIO pin %s", interval, d.pinName)

	initialRead := func() {
		temp, humidity, err := d.Read()
		if err != nil {
			log.Printf("Error reading sensor: %v", err)
			return
		}

		if err := database.InsertMetric(temp, humidity); err != nil {
			log.Printf("Error inserting metric: %v", err)
			return
		}

		log.Printf("Metric saved: Temp=%.1f°C, Humidity=%.1f%%", temp, humidity)
	}

	initialRead()

	for {
		select {
		case <-ctx.Done():
			log.Println("Sensor collector stopped")
			return
		case <-ticker.C:
			initialRead()
		}
	}
}

func (d *DHT11) Close() error {
	return nil
}
