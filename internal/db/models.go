package db

import "time"

type RawMetric struct {
	ID        int64     `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Temp      float64   `json:"temp"`
	Humidity  float64   `json:"humidity"`
}

type DailySummary struct {
	ID          int64     `json:"id"`
	Date        time.Time `json:"date"`
	AvgTemp     float64   `json:"avg_temp"`
	AvgHumidity float64   `json:"avg_humidity"`
	MaxTemp     float64   `json:"max_temp"`
	MinTemp     float64   `json:"min_temp"`
}

type UserLog struct {
	ID         int64     `json:"id"`
	Date       time.Time `json:"date"`
	Rating     int       `json:"rating"`
	Note       string    `json:"note"`
	FeelingTag string    `json:"feeling_tag"`
}
