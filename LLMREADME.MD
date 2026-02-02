# Sistema de Telemetria Térmica e Bem-Estar

Objetivo: Criar uma aplicação em Go que monitore um sensor DHT11, armazene dados no SQLite, gerencie expiração de logs e exponha uma API/Dashboard Web.

## Stack Tecnológica

- **Backend:** Go + Python script (scripts/read_sensor.py)
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
# Sistema de Telemetria Térmica e Bem-Estar

Monitoramento de sensor DHT11 usando uma ponte Go + Python, armazenamento em SQLite3 e dashboard web.

## Stack Tecnológica

- **Backend:** Go (Servidor API e Regras de Negócio)
- **Sensor:** Python 3 + `adafruit-circuitpython-dht` (Ponte de leitura)
- **Frontend:** HTML/JS + Tailwind CSS (CDN) + Chart.js
- **Banco:** SQLite3
- **Deploy:** Binário Go + Scripts Python + systemd (Raspberry Pi)

## Arquitetura de Leitura (Python Wrapper)

Devido a instabilidades em bibliotecas nativas de Go para o DHT11, o sistema utiliza um script Python como ponte:
1. O servidor Go executa `python3 scripts/read_sensor.py <GPIO>`.
2. O script Python utiliza `Adafruit_Blinka` para acessar os pinos e retorna um JSON.
3. O Go processa o JSON e salva no banco de dados.

**Pino Padrão:** GPIO17 (Físico 11)

## Estrutura do Projeto

```
raspberryTemperatureSensor/
├── cmd/server/main.go      # Ponto de entrada (Servidor HTTP)
├── internal/
│   ├── api/handler.go      # Endpoints REST
│   ├── config/config.go    # Configurações (GPIO, Intervalos)
│   ├── db/database.go      # Persistência SQLite3
│   ├── sensor/dht11.go     # Lógica de chamada da ponte Python
│   └── maintenance/        # Worker de limpeza de logs antigos
├── scripts/
│   ├── read_sensor.py      # Script Python (leitura real do sensor)
│   └── build-gpio17.sh     # Script oficial de build e deploy
├── static/index.html       # Dashboard Frontend
└── data/telemetry.db       # Banco de dados (gerado automaticamente)
```

## Guia de Operação (Manual)

### 1. Requisitos no Raspberry Pi
É necessário instalar as dependências de hardware para o Python:
```bash
sudo pip3 install Adafruit-Blinka adafruit-circuitpython-dht --break-system-packages
```

### 2. Build e Deploy
Sempre use o script oficial para garantir que todas as partes (Go e Python) sejam movidas para `/opt/mysensoringo`:
```bash
./scripts/build-gpio17.sh
```

### 3. Verificação
- **Logs do Serviço:** `sudo journalctl -u mysensoringo -f`
- **Dashboard:** Acesse `http://<ip-do-raspberry>:8080`

## Endpoints Principais
- `GET /api/current`: Última temperatura e umidade.
- `GET /api/history`: Histórico completo para os gráficos.
- `POST /api/feeling`: Salvar sensação térmica subjetiva (1-5 estrelas).

## Manutenção Automática
O sistema possui um worker que roda a cada 24h para:
1. Consolidar médias diárias em `daily_summaries`.
2. Remover dados brutos (`raw_metrics`) com mais de 30 dias.
└── DEPLOY.md
