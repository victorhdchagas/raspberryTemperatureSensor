package web

import (
	"fmt"
	"time"
)

const (
	templateCurrentReading = `
		<div id="currentReading" class="text-center">
			<p class="text-4xl font-bold text-green-400 mb-2">%s°C</p>
			<p class="text-2xl text-blue-400">%s%%</p>
			<p class="text-sm text-gray-400 mt-2">%s</p>
		</div>
	`

	templateNoData = `
		<p class="text-gray-400">Sem dados disponíveis</p>
	`

	templateHotDay = `
		<div class="bg-gray-700 rounded p-4">
			<p class="font-semibold text-yellow-400">%s</p>
			<p class="text-2xl font-bold">%.1f°C</p>
			<p class="text-sm text-gray-400">Média</p>
			<p class="text-sm text-gray-400">Max: %.1f°C | Min: %.1f°C</p>
		</div>
	`

	templateError = `
		<p class="text-red-400">Erro ao carregar dados</p>
	`
)

func RenderCurrentReading(temp, humidity float64, timestamp time.Time) string {
	return fmt.Sprintf(templateCurrentReading,
		fmt.Sprintf("%.1f", temp),
		fmt.Sprintf("%.1f", humidity),
		timestamp.Format("02/01/2006 15:04:05"),
	)
}

func RenderNoData() string {
	return templateNoData
}

func RenderHotDay(date string, avgTemp, maxTemp, minTemp float64) string {
	return fmt.Sprintf(templateHotDay, date, avgTemp, maxTemp, minTemp)
}

func RenderError() string {
	return templateError
}
