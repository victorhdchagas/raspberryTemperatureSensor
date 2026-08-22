package weather

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// cacheTTL controla a frequência máxima de consultas à API. Com TTL de 1h,
// o serviço faz no máximo ~24 consultas/dia, longe de qualquer rate limit.
const cacheTTL = time.Hour

// Client consulta a temperatura externa de uma região via Open-Meteo,
// com cache em memória para evitar rate limiting da API.
type Client struct {
	Lat    float64
	Lon    float64
	mu     sync.Mutex
	cached *cachedReading
	http   *http.Client
}

type cachedReading struct {
	Temp      float64
	FetchedAt time.Time
}

func New(lat, lon float64) *Client {
	return &Client{
		Lat:  lat,
		Lon:  lon,
		http: &http.Client{Timeout: 10 * time.Second},
	}
}

// GetTemp retorna a temperatura externa atual, usando cache se a última
// consulta foi feita há menos de cacheTTL. Se falhar, retorna o último valor
// conhecido (stale) para não quebrar a leitura do sensor.
func (c *Client) GetTemp() (*float64, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cached != nil && time.Since(c.cached.FetchedAt) < cacheTTL {
		return &c.cached.Temp, nil
	}

	temp, err := c.fetch()
	if err != nil {
		if c.cached != nil {
			log.Printf("weather: fetch falhou, usando cache stale (%.1f°C): %v", c.cached.Temp, err)
			return &c.cached.Temp, nil
		}
		return nil, err
	}

	c.cached = &cachedReading{Temp: temp, FetchedAt: time.Now()}
	return &temp, nil
}

type openMeteoResponse struct {
	Current struct {
		Temperature2M float64 `json:"temperature_2m"`
	} `json:"current"`
}

func (c *Client) fetch() (float64, error) {
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.4f&longitude=%.4f&current=temperature_2m&timezone=America%%2FSao_Paulo",
		c.Lat, c.Lon,
	)

	resp, err := c.http.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("open-meteo HTTP %d", resp.StatusCode)
	}

	var data openMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return 0, err
	}

	return data.Current.Temperature2M, nil
}
