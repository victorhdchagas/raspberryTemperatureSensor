# MySensorIngo - Sistema de Telemetria Térmica

Sistema em Go para monitoramento de temperatura e umidade com sensor DHT11 no Raspberry Pi.

## 🚀 Deploy Rápido (RESOLVIDO GPIO 4)

### ⚠️ PROBLEMA: GPIO 4 não funciona!

**Razão:** GPIO 4 conflita com sistema 1-Wise do Raspberry Pi (conflito de hardware).

### ✅ SOLUÇÃO RÁPIDA (Recomendada)

Execute este comando no Raspberry Pi:

```bash
cd /opt/mysensoringo
bash scripts/build-gpio17.sh
```

**E MUDAR O CONECTOR FÍSICO:**
- Mova o fio **DATA** do DHT11
- De: **Pin 7** (GPIO 4)
- Para: **Pin 11** (GPIO 17)

```
DHT11       Raspberry Pi
---------   ---------------
VCC   -----> Pin 1  (3.3V)
GND   -----> Pin 6  (GND)
DATA  -----> Pin 11 (GPIO 17) ← MUDAR!
```

### Outras Opções

#### Opção 2: Deploy Normal (GPIO 4 - Pode não funcionar)

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

#### Opção 3: Deploy Remoto (via SSH)

```bash
bash scripts/deploy.sh
```

#### Opção 4: Testar Pinos Automaticamente

```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
```

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

### Pinos GPIO (IMPORTANTE!)

| GPIO | Physical | Status | Uso |
|-------|----------|---------|------|
| 4 | Pin 7 | ❌ **Conflita 1-Wise** - Não usar |
| 17 | Pin 11 | ✅ **Recomendado** - Use este! |
| 18 | Pin 12 | ✅ Funciona (mas pode conflitar áudio) |
| 27 | Pin 13 | ✅ Funciona |
| 22 | Pin 15 | ✅ Funciona |

### Mudar Pino GPIO

#### Opção 1: Via script (Recomendado)
```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh  # Testa e muda automaticamente
# OU
bash scripts/build-gpio17.sh  # Usa GPIO 17 diretamente
```

#### Opção 2: Via código
Editar `internal/config/config.go`:
```go
Sensor: SensorConfig{
    GPIO:     "17",  // MUDAR DE "4" PARA "17"
    Interval: 50 * time.Minute,
},
```

#### Opção 3: Via variável de ambiente
```bash
export SENSOR_GPIO_PIN=17
cd /opt/mysensoringo
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

## 📌 Notas GPIO Importantes

### Por que GPIO 4 não funciona?

O **GPIO 4 (BCM)** é usado automaticamente pelo **sistema 1-Wise** do Raspberry Pi para sensores de temperatura DS18B20. Isso é um conflito de hardware que não pode ser resolvido apenas com permissões.

### Solução: Usar GPIO 17

GPIO 17 (physical pin 11) é seguro e não conflita com nenhum sistema do Raspberry Pi.

### Conexão DHT11 com GPIO 17

```
DHT11       Raspberry Pi
---------   ---------------
VCC   -----> Pin 1  (3.3V)
GND   -----> Pin 6  (GND)
DATA  -----> Pin 11 (GPIO 17) ← MUDAR DE PIN 7
```

## 📝 Nota Importante

O sensor faz a **primeira leitura imediatamente** ao iniciar, seguida por leituras a cada 50 minutos (configurável).
