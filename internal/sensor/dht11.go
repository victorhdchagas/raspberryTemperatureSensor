package sensor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
	"github.com/wutachi/raspberryTemperatureSensor/internal/weather"
)

type Reading struct {
	Temp     float64 `json:"temperature"`
	Humidity float64 `json:"humidity"`
	Error    *string `json:"error"`
}

type DHT11 struct {
	pinName string
	weather *weather.Client
}

func NewDHT11(pinName string, weatherClient *weather.Client) (*DHT11, error) {
	return &DHT11{
		pinName: pinName,
		weather: weatherClient,
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

		var tempExt *float64
		if d.weather != nil {
			if ext, wErr := d.weather.GetTemp(); wErr != nil {
				log.Printf("Warning: could not fetch external temp: %v", wErr)
			} else {
				tempExt = ext
			}
		}

		if err := database.InsertMetric(temp, humidity, tempExt); err != nil {
			log.Printf("Error inserting metric: %v", err)
			return
		}

		if tempExt != nil {
			log.Printf("Metric saved: Temp=%.1f°C, Humidity=%.1f%%, Ext=%.1f°C", temp, humidity, *tempExt)
		} else {
			log.Printf("Metric saved: Temp=%.1f°C, Humidity=%.1f%%", temp, humidity)
		}
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

func (d *DHT11) StartAsync(ctx context.Context, interval time.Duration, database *db.Database) {
	go d.Start(ctx, interval, database)
}

func (d *DHT11) Close() error {
	return nil
}
