# MySensorIngo - Sistema de Telemetria Térmica

Sistema em Go para monitoramento de temperatura e umidade com sensor DHT11 no Raspberry Pi.

## 🚀 Deploy Rápido

### Opção 1: Deploy Local (Recomendado)

Execute diretamente no Raspberry Pi:

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

✅ **Primeira leitura imediata** ao iniciar o serviço!
⏹️ **Para o serviço** antes de copiar binário (evita erro "Text file busy")
⚠️ **Importante:** O serviço roda como `root` para garantir acesso ao GPIO.

### Opção 2: Deploy Remoto (via SSH)

```bash
bash scripts/deploy.sh
```

## 📋 Funcionalidades

- 🌡️ **Leitura de sensor:** DHT11 via GPIO (primeira leitura imediata + 50min)
- 📊 **Dashboard web:** Gráficos, histórico e calendário térmico
- 🗄️ **SQLite3:** Armazenamento local com retenção de 30 dias
- 🧹 **Auto-limpeza:** Worker que deleta dados antigos diariamente
- 📝 **Notas pessoais:** Sistema de rating (1-5 estrelas) e observações

## 📡 API Endpoints

| Endpoint | Método | Descrição |
|----------|---------|-------------|
| `/api/current` | GET | Última leitura do sensor |
| `/api/history` | GET | Histórico por período |
| `/api/stats/hot-days` | GET | Dias mais quentes |
| `/api/feeling` | POST | Salvar sensação do dia |

## 🔧 Configuração

### Opção 1: Via código
Editar `internal/config/config.go` para alterar:
- Pino GPIO: padrão = 4 (BCM)
- Intervalo de leitura: padrão = 50 minutos
- Porta HTTP: padrão = 8080
- Path do banco: padrão = `./data/telemetry.db`

### Opção 2: Via variável de ambiente
```bash
export SENSOR_GPIO_PIN=17  # Muda pino para GPIO 17
bash scripts/build-local.sh
```

## 📱 Dashboard

Acesse em:
```
http://<ip-do-raspberry>:8080
```

## 🛠️ Instalação Completa

Veja `DEPLOY_LOCAL.md` para instruções detalhadas.

## ⚠️ Troubleshooting

Se encontrar erro "failed to export pin 4":
- Veja `TROUBLESHOOTING.md` para soluções completas
- O serviço roda como `root` para evitar problemas de permissão
- Se necessário, mude para outro pino (ex: GPIO 17)

## 📝 Documentação Adicional

- `DEPLOY_LOCAL.md` - Guia de instalação local
- `TROUBLESHOOTING.md` - Solução de problemas GPIO
- `UPDATE.md` - Atualizações importantes

## 🏗️ Desenvolvimento

```bash
# Instalar dependências
go mod download

# Build local (para testes)
go build -o telemetry-server ./cmd/server
./telemetry-server

# Cross-compile (complexo, requer GCC ARM)
GOOS=linux GOARCH=arm GOARM=7 go build -o telemetry-server ./cmd/server
```

## 📁 Estrutura

```
├── cmd/server/         # Entry point
├── internal/
│   ├── config/         # Configurações
│   ├── db/           # SQLite3 + migrations
│   ├── sensor/       # DHT11 reader (primeira leitura imediata)
│   ├── api/          # Endpoints REST
│   └── maintenance/   # Worker de limpeza
├── pkg/web/          # Arquivos estáticos
├── static/           # Dashboard (HTML/JS + Tailwind + Chart.js)
└── scripts/
    ├── build.sh       # Build cross-platform
    ├── build-local.sh # Build + deploy local
    └── deploy.sh     # Deploy remoto via SSH
```

## ⚙️ Requisitos

- Go 1.21+
- GCC (para CGO)
- Raspberry Pi com DHT11
- systemd
- Permissões de root (serviço roda como root)

## 📌 Notas GPIO

- GPIO 4 (BCM) é padrão, mas pode conflitar com 1-Wire
- Recomenda GPIO 17 (physical pin 11) se houver problemas
- O serviço roda como `root` para garantir acesso ao GPIO

## 📝 Nota Importante

O sensor faz a **primeira leitura imediatamente** ao iniciar, seguida por leituras a cada 50 minutos (configurável).
