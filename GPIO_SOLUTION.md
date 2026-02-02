# Solução Definitiva - GPIO 4 e 1-Wire

## 🎯 O Problema

O erro **"failed to export pin 4"** persiste porque:

1. GPIO 4 (BCM numbering) é usado por padrão pelo **sistema 1-Wire** do Raspberry Pi
2. Rodar como `root` não resolve - é um conflito de hardware
3. O 1-Wise monopoliza o pino, impedindo acesso direto

---

## ✅ SOLUÇÃO: Mudar para GPIO 17 (Recomendado)

GPIO 17 é seguro e não conflita com nenhum sistema do Raspberry Pi.

### Método 1: Usar script automático (Recomendado)

```bash
# No Raspberry Pi
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
```

O script vai:
1. Testar pinos disponíveis (4, 17, 18, 27, 22)
2. Mostrar quais funcionam
3. Mudar automaticamente para um pino que funciona

### Método 2: Mudar manualmente

1. Editar config:
```bash
nano /opt/mysensoringo/internal/config/config.go
```

2. Mudar:
```go
Sensor: SensorConfig{
    GPIO:     "17",  // MUDAR DE "4" PARA "17"
    Interval: 50 * time.Minute,
},
```

3. Rebuild:
```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## 🔌 Conexão Física

Se mudar para GPIO 17:

| GPIO 17 (BCM) | Physical Pin | Descrição |
|----------------|---------------|-------------|
| GPIO 17 | Pin 11 | Data do DHT11 ✅ |

Conectar o DHT11:
- VCC → 3.3V (Pin 1 ou 17)
- GND → GND (Pin 6, 9, 14, 20, 25, 30, 34 ou 39)
- Data → GPIO 17 (Pin 11)

---

## ⚙️ Outros Pinos Testados

| GPIO | Physical | Status |
|-------|----------|---------|
| 4 | Pin 7 | ❌ Conflita com 1-Wire |
| 17 | Pin 11 | ✅ **Recomendado** |
| 18 | Pin 12 | ✅ Funciona (mas pode conflitar com áudio) |
| 27 | Pin 13 | ✅ Funciona |
| 22 | Pin 15 | ✅ Funciona |

---

## 🔧 Desabilitar 1-Wire (Alternativa)

Se quiser MANTER GPIO 4, precisa desabilitar 1-Wise:

```bash
sudo nano /boot/config.txt
```

Adicionar/modificar:
```
dtoverlay=w1-gpio-pullup,gpiopin=18
```

Isso move 1-Wire para GPIO 18, liberando GPIO 4. Depois:
```bash
sudo reboot
```

---

## 📋 Resumo

**Recomendação:** Use GPIO 17 (physical pin 11) - é o mais seguro e simples!

```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
# Ou manualmente:
nano internal/config/config.go  # Mude GPIO: "4" para GPIO: "17"
bash scripts/build-local.sh
```

---

## ✅ Verificar se Funcionou

```bash
sudo journalctl -u mysensoringo -f
```

Deve ver:
```
Metric saved: Temp=XX.X°C, Humidity=XX.X%
```

**Se ainda der erro:** Use o script `fix-gpio.sh` para testar pinos diferentes automaticamente!
