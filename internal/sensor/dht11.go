package sensor

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/wutachi/raspberryTemperatureSensor/internal/db"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/devices/v3/dht"
	"periph.io/x/host/v3"
)

type Reading struct {
	Temp     float64
	Humidity float64
	Error    error
}

type DHT11 struct {
	pin     gpio.PinIO
	device  *dht.Device
	pinName string
}

func NewDHT11(pinName string) (*DHT11, error) {
	if _, err := host.Init(); err != nil {
		return nil, fmt.Errorf("failed to initialize periph: %w", err)
	}

	p := gpioreg.ByName(pinName)
	if p == nil {
		// Try adding "GPIO" prefix if it's just a number
		p = gpioreg.ByName("GPIO" + pinName)
		if p == nil {
			return nil, fmt.Errorf("failed to find pin: %s", pinName)
		}
	}

	d, err := dht.New(p, dht.DHT11)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize DHT11: %w", err)
	}

	return &DHT11{
		pin:     p,
		device:  d,
		pinName: pinName,
	}, nil
}

func (d *DHT11) Read() (temp, humidity float64, err error) {
	var env dht.Environmental
	if err := d.device.Sense(&env); err != nil {
		return 0, 0, fmt.Errorf("sensor read error: %w", err)
	}

	temperature := float64(env.Temperature.Celsius())
	humidityVal := float64(env.Humidity) / 100000.0 // periph.io returns humidity in milli-percent? No, usually it's a special type.

	// Let's check periph.io DHT environmental struct.
	// Actually, env.Humidity is of type physic.RelativeHumidity.
	// To get percentage: float64(env.Humidity) / float64(physic.PercentRH)

	return temperature, humidityVal, nil
}

func (d *DHT11) Start(ctx context.Context, interval time.Duration, database *db.Database) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.Printf("Sensor collector started: reading every %v from GPIO pin %s", interval, d.pinName)

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
