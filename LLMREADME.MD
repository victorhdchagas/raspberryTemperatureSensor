# Sistema de Telemetria Térmica e Bem-Estar

Objetivo: Criar uma aplicação em Go que monitore um sensor DHT11, armazene dados no SQLite, gerencie expiração de logs e exponha uma API/Dashboard Web.

## Stack Tecnológica

- **Backend:** Go + d2r2/go-dht (GPIO)
- **Frontend:** HTML/JS + Tailwind CSS (CDN) + Chart.js
- **Banco:** SQLite3
- **Deploy:** Cross-compile ARM7 + SCP + systemd

## Arquitetura do Projeto

### 1. Arquitetura de Dados (Database & Migrations)

- **Engine:** SQLite3
- **Migrations:** golang-migrate ou script interno
- **Tabelas:**
  - `raw_metrics`: (id, timestamp, temp, humidity)
  - `daily_summaries`: (id, date, avg_temp, avg_humidity, max_temp, min_temp)
  - `user_logs`: (id, date, rating [1-5], note, feeling_tag)
- **Rotina de Maintenance:** Worker em background (cada 24h)
  - Calcula média do dia anterior → `daily_summaries`
  - Deleta `raw_metrics` com mais de 30 dias

### 2. Backend (Golang API)

- **Coletor:** Goroutine que lê o sensor (GPIO) a cada 50 min com tratamento de erro "buffer not returned"
- **Endpoints:**
  - `GET /api/current`: Retorna a última leitura e status do sensor
  - `GET /api/history?start=...&end=...`: Filtro por período
  - `GET /api/stats/hot-days`: Dias com maiores médias
  - `POST /api/feeling`: Salva rating e nota do dia (sistema de sensação)

### 3. Frontend (Web Dashboard)

- **Estilo:** Minimalista, "clean interface"
- **Componentes:**
  - Calendário Térmico: GitHub Contribution Graph style (cor varia por temperatura/rating)
  - Filtros: Seletor de data com gráficos de linha (Chart.js)
  - Modal de Notas: Adicionar/editar nota e rating (1-5 estrelas)

### 4. Deploy (Raspberry Pi 3)

- **Cross-compile:** GOOS=linux GOARCH=arm GOARM=7
- **Script de deploy:** scp + reiniciar serviço systemd

## Checklist de Tarefas

### 📦 Fase 1 - Infraestrutura (Tasks 1-3)

- [x] Task 1: Configurar estrutura do projeto Go (go mod, diretórios cmd, internal, pkg, web)
- [x] Task 2: Implementar schema do SQLite3 e sistema de migrations (script interno)
- [x] Task 3: Criar tabelas: raw_metrics, daily_summaries, user_logs

### 🌡️ Fase 2 - Sensor & Coleta (Tasks 4-5)

- [x] Task 4: Implementar leitura do sensor DHT11 via GPIO usando d2r2/go-dht
- [x] Task 5: Criar goroutine coletora que lê sensor a cada 50 min com tratamento de erro

### 🔌 Fase 3 - Backend API (Tasks 6-9)

- [x] Task 6: Implementar API GET /api/current (última leitura + status)
- [x] Task 7: Implementar API GET /api/history (filtro por período)
- [x] Task 8: Implementar API GET /api/stats/hot-days (dias com maiores médias)
- [x] Task 9: Implementar API POST /api/feeling (salvar rating e nota)

### 🧹 Fase 4 - Maintenance (Tasks 10-12)

- [x] Task 10: Criar worker background de maintenance (rodar a cada 24h)
- [x] Task 11: Implementar lógica de cálculo de médias diárias e inserção em daily_summaries
- [x] Task 12: Implementar deleção de raw_metrics com mais de 30 dias

### 🎨 Fase 5 - Frontend (Tasks 13-16)

- [x] Task 13: Criar frontend HTML/JS com Tailwind CSS (CDN)
- [x] Task 14: Implementar Calendário Térmico (GitHub contribution graph style)
- [x] Task 15: Adicionar filtros de data e gráficos de linha (Chart.js)
- [x] Task 16: Implementar modal de notas para adicionar/editar rating (1-5 estrelas)

### 🚀 Fase 6 - Build & Deploy (Tasks 17-18)

- [x] Task 17: Configurar cross-compile para GOOS=linux GOARCH=arm GOARM=7
- [x] Task 18: Criar script de deploy (scp + reiniciar serviço systemd)

## Sensor Setup

- **Pino GPIO:** D4 (conforme teste_sensor.py)
- **Sensor:** DHT11
- **Biblioteca Python (referência):** adafruit_dht

## Status Atual

**Progresso:** 18/18 tarefas concluídas ✅

**Última atualização:** 2026-02-01

## Estrutura do Projeto

```
raspberryTemperatureSensor/
├── cmd/
│   └── server/
│       └── main.go              # Ponto de entrada do servidor
├── internal/
│   ├── config/
│   │   └── config.go           # Configurações da aplicação
│   ├── db/
│   │   └── database.go         # SQLite3 e migrations
│   ├── sensor/
│   │   └── dht11.go           # Leitura do sensor DHT11
│   ├── api/
│   │   └── handler.go         # Endpoints da API REST
│   └── maintenance/
│       └── worker.go          # Worker de manutenção
├── pkg/
│   └── web/
│       └── static.go          # Servidor de arquivos estáticos
├── static/
│   └── index.html             # Dashboard web
├── scripts/
│   ├── build.sh               # Script de build (cross-compile)
│   └── deploy.sh              # Script de deploy (scp + systemd)
├── configs/
│   ├── telemetry-service.service  # Configuração systemd
│   └── deploy.env              # Variáveis de ambiente do deploy
├── build/
│   └── telemetry-server       # Binário compilado
├── go.mod
├── go.sum
├── LLMREADME.MD
└── DEPLOY.md
```
