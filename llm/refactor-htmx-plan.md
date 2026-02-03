# Plano de Refatoração e Implementação HTMX

## Visão Geral

Este plano detalha a refatoração do código atual para seguir princípios de Clean Code, modularização e preparação para implementação de HTMX como frontend framework.

## Análise Atual

### Estrutura do Projeto
```
raspberryTemperatureSensor/
├── cmd/
│   └── server/main.go (79 linhas) ✅
├── internal/
│   ├── api/
│   │   └── handler.go (179 linhas) ❌ EXCEDE LIMITE
│   ├── config/
│   │   └── config.go (49 linhas) ✅
│   ├── db/
│   │   └── database.go (229 linhas) ❌ EXCEDE LIMITE
│   ├── maintenance/
│   │   └── worker.go (68 linhas) ✅
│   └── sensor/
│       └── dht11.go (87 linhas) ✅
├── pkg/
│   └── web/static.go (11 linhas) ✅
└── static/
    ├── index.html (15KB)
    ├── css/tailwind.css (8.3KB)
    └── js/chart.min.js (204KB) ❌ MUITO GRANDE
```

### Problemas Identificados

#### 1. Falta de Modularização
- `internal/db/database.go` (229 linhas) deve ser dividido
- `internal/api/handler.go` (179 linhas) deve ser dividido

#### 2. Dependência Excessiva de JSON
- Todos os endpoints retornam apenas JSON
- Sem suporte a renderização HTML (necessário para HTMX)
- Sem Server-Sent Events (SSE) para atualizações em tempo real

#### 3. Frontend Pesado
- Chart.js completo (204KB) desnecessário para o Raspberry Pi 3
- JavaScript inline no HTML (330 linhas misturadas com HTML)
- Polling a cada 60s (ineficiente)

#### 4. Clean Code Violations
- Constantes espalhadas em vários lugares
- Funções com múltiplas responsabilidades
- Falta de separação de conceitos

## Plano de Refatoração (Clean Code + Modularização)

### Fase 1: Refatoração do Banco de Dados (db/)

#### Arquivo: `internal/db/models.go` (novo)
**Propósito:** Definição de todos os modelos de dados
**Tamanho estimado:** ~40 linhas

```go
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
```

#### Arquivo: `internal/db/database.go` (refatorado)
**Propósito:** Core do database, conexão e operações básicas
**Tamanho estimado:** ~60 linhas

- Mantém: `New()`, `Close()`, Ping básico
- Remove: Todos os models, Migrate, métodos específicos

#### Arquivo: `internal/db/migrations.go` (novo)
**Propósito:** Definição das migrations
**Tamanho estimado:** ~50 linhas

- Extrai todas as queries SQL de criação de tabelas
- Função `GetMigrations()` que retorna slice de strings

#### Arquivo: `internal/db/metrics.go` (novo)
**Propósito:** Operações relacionadas a raw_metrics
**Tamanho estimado:** ~60 linhas

- `InsertMetric(temp, humidity float64)`
- `GetLatestMetric()`
- `GetMetricsByDateRange(start, end)`
- `DeleteOldMetrics(days)`

#### Arquivo: `internal/db/summaries.go` (novo)
**Propósito:** Operações relacionadas a daily_summaries
**Tamanho estimado:** ~50 linhas

- `InsertDailySummary(summary)`
- `GetHotDays(limit)`
- `CalculateDailySummary(date)`

#### Arquivo: `internal/db/userlogs.go` (novo)
**Propósito:** Operações relacionadas a user_logs
**Tamanho estimado:** ~40 linhas

- `InsertUserLog(log)`
- `GetUserLogByDate(date)`

#### Arquivo: `internal/db/constants.go` (novo)
**Propósito:** Constantes relacionadas ao banco de dados
**Tamanho estimado:** ~20 linhas

```go
package db

const (
    DefaultRetentionDays = 30
    DefaultHotDaysLimit  = 10
    // Outras constantes...
)
```

### Fase 2: Refatoração da API (api/)

#### Arquivo: `internal/api/handler.go` (refatorado)
**Propósito:** Core do handler, registro de rotas
**Tamanho estimado:** ~50 linhas

- Mantém: `NewHandler`, `RegisterRoutes`
- Remove: Todos os handlers específicos, funções utilitárias

#### Arquivo: `internal/api/current.go` (novo)
**Propósito:** Handler para endpoint `/api/current`
**Tamanho estimado:** ~60 linhas

- `func (h *Handler) getCurrent(w, r)`
- `func (h *Handler) getCurrentHTML(w, r)` (para HTMX)

#### Arquivo: `internal/api/history.go` (novo)
**Propósito:** Handler para endpoint `/api/history`
**Tamanho estimado:** ~60 linhas

- `func (h *Handler) getHistory(w, r)`
- Validação de parâmetros

#### Arquivo: `internal/api/hotdays.go` (novo)
**Propósito:** Handler para endpoint `/api/stats/hot-days`
**Tamanho estimado:** ~50 linhas

- `func (h *Handler) getHotDays(w, r)`
- Validação de limit

#### Arquivo: `internal/api/feeling.go` (novo)
**Propósito:** Handler para endpoint `/api/feeling`
**Tamanho estimado:** ~70 linhas

- `func (h *Handler) postFeeling(w, r)`
- Validação de request

#### Arquivo: `internal/api/response.go` (novo)
**Propósito:** Funções utilitárias de resposta HTTP
**Tamanho estimado:** ~40 linhas

- `func respondJSON(w, status, data)`
- `func respondHTML(w, status, html)`
- `func respondError(w, status, message)`

#### Arquivo: `internal/api/validation.go` (novo)
**Propósito:** Funções de validação
**Tamanho estimado:** ~40 linhas

- `func parseLimit(s string)`
- Outras validações comuns

#### Arquivo: `internal/api/sse.go` (novo)
**Propósito:** Server-Sent Events para streaming em tempo real
**Tamanho estimado:** ~80 linhas

- `func (h *Handler) streamCurrent(w, r)` (endpoint `/api/current/stream`)
- Gerenciamento de conexões SSE

### Fase 3: Refatoração do Web Server (web/)

#### Arquivo: `pkg/web/templates.go` (novo)
**Propósito:** Templates HTML para HTMX
**Tamanho estimado:** ~100 linhas

- Templates inline ou carregados de arquivos
- Template para componente de temperatura atual
- Template para componente de calendário térmico
- Template para componente de dias quentes

#### Arquivo: `pkg/web/static.go` (refatorado)
**Propósito:** Servir arquivos estáticos
**Tamanho estimado:** ~30 linhas

- Mantém: Registro de rotas estáticas
- Adiciona: Headers de cache

## Plano de Implementação HTMX

### Fase 4: Implementação HTMX no Frontend

#### Arquivo: `static/index.html` (refatorado)
**Principais mudanças:**
- Remove JavaScript inline (330 linhas)
- Adiciona HTMX via CDN (~14KB)
- Substitui Chart.js completo por Chart.js Lite (~40KB)
- Usa atributos HTMX para atualizações

**Estrutura do novo index.html:**
```html
<!DOCTYPE html>
<html lang="pt-BR">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Telemetria Térmica</title>
    <link rel="stylesheet" href="/static/css/tailwind.css">
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <script src="/static/js/chart-lite.min.js"></script>
</head>
<body class="bg-gray-900 text-white">
    <!-- Componentes HTMX -->
    <div id="current"
         hx-get="/api/current/html"
         hx-trigger="load, every 60s"
         hx-swap="outerHTML">
        Carregando...
    </div>
    
    <!-- Outros componentes -->
</body>
</html>
```

#### Arquivo: `static/css/htmx.css` (novo)
**Propósito:** Estilos específicos para HTMX
**Tamanho estimado:** ~30 linhas

- Classes para transições
- Estilos para loading states

### Fase 5: Implementação SSE (Server-Sent Events)

#### Endpoint: `/api/current/stream`
**Propósito:** Streaming de dados em tempo real
**Arquivo:** `internal/api/sse.go`

**Benefícios:**
- Substitui polling (ineficiente)
- Atualizações instantâneas quando há novos dados
- Menor tráfego de rede
- Menor carga no servidor

## Estrutura Final do Projeto

```
raspberryTemperatureSensor/
├── cmd/
│   └── server/main.go (79 linhas) ✅
├── internal/
│   ├── api/
│   │   ├── handler.go (~50 linhas) ✅
│   │   ├── current.go (~60 linhas) ✅
│   │   ├── history.go (~60 linhas) ✅
│   │   ├── hotdays.go (~50 linhas) ✅
│   │   ├── feeling.go (~70 linhas) ✅
│   │   ├── response.go (~40 linhas) ✅
│   │   ├── validation.go (~40 linhas) ✅
│   │   └── sse.go (~80 linhas) ✅
│   ├── config/
│   │   └── config.go (49 linhas) ✅
│   ├── db/
│   │   ├── database.go (~60 linhas) ✅
│   │   ├── models.go (~40 linhas) ✅
│   │   ├── migrations.go (~50 linhas) ✅
│   │   ├── metrics.go (~60 linhas) ✅
│   │   ├── summaries.go (~50 linhas) ✅
│   │   ├── userlogs.go (~40 linhas) ✅
│   │   └── constants.go (~20 linhas) ✅
│   ├── maintenance/
│   │   └── worker.go (68 linhas) ✅
│   └── sensor/
│       └── dht11.go (87 linhas) ✅
├── pkg/
│   ├── web/
│   │   ├── static.go (~30 linhas) ✅
│   │   └── templates.go (~100 linhas) ✅
└── static/
    ├── index.html (refatorado, ~80 linhas) ✅
    ├── css/
    │   ├── tailwind.css (8.3KB)
    │   └── htmx.css (~30 linhas) ✅
    └── js/
        └── chart-lite.min.js (~40KB, substitui chart.min.js 204KB) ✅
```

## Passos de Implementação

### Passo 1: Refatoração do Banco de Dados (Sem quebrar funcionalidade)

1. Criar `internal/db/models.go` com todos os modelos
2. Criar `internal/db/constants.go` com constantes
3. Criar `internal/db/metrics.go` e mover métodos relacionados
4. Criar `internal/db/summaries.go` e mover métodos relacionados
5. Criar `internal/db/userlogs.go` e mover métodos relacionados
6. Criar `internal/db/migrations.go` e extrair migrations
7. Refatorar `internal/db/database.go` para core operations
8. Testar que tudo ainda funciona

### Passo 2: Refatoração da API (Sem quebrar funcionalidade)

1. Criar `internal/api/response.go` com funções utilitárias
2. Criar `internal/api/validation.go` com validações
3. Criar `internal/api/current.go` e mover handler
4. Criar `internal/api/history.go` e mover handler
5. Criar `internal/api/hotdays.go` e mover handler
6. Criar `internal/api/feeling.go` e mover handler
7. Refatorar `internal/api/handler.go` para core
8. Testar que tudo ainda funciona

### Passo 3: Preparação para HTMX (Endpoints HTML)

1. Adicionar métodos HTML em `internal/api/current.go`
2. Adicionar métodos HTML em `internal/api/hotdays.go`
3. Criar `internal/api/sse.go` com SSE endpoint
4. Criar `pkg/web/templates.go` com templates HTML
5. Testar endpoints HTML e SSE

### Passo 4: Implementação Frontend HTMX

1. Baixar Chart.js Lite (~40KB) e substituir Chart.js completo
2. Criar `static/css/htmx.css`
3. Refatorar `static/index.html` com HTMX
4. Remover JavaScript inline
5. Adicionar HTMX via CDN
6. Testar que tudo funciona

### Passo 5: Testes e Validação

1. Testar todos os endpoints existentes (backward compatibility)
2. Testar novos endpoints HTML
3. Testar SSE streaming
4. Medir performance (antes e depois)
5. Testar no Raspberry Pi 3

## Benefícios Esperados

### Performance
- **Redução de 85% no tamanho do JavaScript** (204KB → 30KB)
- **Carregamento inicial instantâneo** (HTML estático)
- **Atualizações em tempo real sem polling** (SSE)
- **Menos tráfego de rede**

### Manutenibilidade
- **Arquivos pequenos e focados** (todos < 100 linhas)
- **Separação clara de responsabilidades**
- **Clean code principles aplicados**
- **Fácil adicionar novos endpoints**

### Recursos do Raspberry Pi
- **Menor uso de RAM** (Chart.js 204KB removido)
- **Menor uso de CPU** (sem polling)
- **Mais espaço em disco** (frontend mais leve)

## Cronograma Estimado

- **Passo 1:** 30-45 minutos
- **Passo 2:** 30-45 minutos
- **Passo 3:** 45-60 minutos
- **Passo 4:** 60-90 minutos
- **Passo 5:** 30-45 minutos

**Total:** ~3-4 horas

## Notas Importantes

1. **NÃO REMOVER ARQUIVOS EXISTENTES** - apenas criar novos e refatorar
2. **MANTER BACKWARD COMPATIBILITY** - endpoints JSON continuam funcionando
3. **TESTAR A CADA PASSO** - garantir que nada quebra
4. **HTMX É ADICIONAL** - não é obrigatório usar imediatamente
5. **SSE É OPCIONAL** - pode ser implementado depois se necessário

## Próximos Passos

Após aprovação deste plano:
1. Começar pela Passo 1 (Refatoração DB)
2. Seguir passo a passo
3. Documentar cada mudança
4. Commitar em etapas
5. Deploy e testes
