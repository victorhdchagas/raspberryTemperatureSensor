# 🚨 SOLUÇÃO RÁPIDA - Erro GPIO 4

## O Problema

**GPIO 4 CONFLITA com 1-Wise** - rodar como `root` não resolve!

Isso é um conflito de hardware no Raspberry Pi. O GPIO 4 é usado automaticamente pelo sistema 1-Wire.

---

## ✅ SOLUÇÃO IMEDIATA (Recomendada)

Execute este comando no Raspberry Pi:

```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
```

### O que vai acontecer:

1. 🔍 **Testar pinos** - Vai testar GPIO 4, 17, 18, 27, 22
2. 📋 **Mostrar resultados** - Quais funcionam, quais dão erro
3. ⚙️ **Mudar automaticamente** - Configura um pino que funciona

---

## 🔌 MUDAR O CONECTOR FÍSICO

Se o script mudar para GPIO 17 (recomendado):

**Mova o fio DATA do DHT11:**
- De: **Pin 7** (GPIO 4)
- Para: **Pin 11** (GPIO 17)

### Diagrama:

```
DHT11       Raspberry Pi
---------   ---------------
VCC   -----> Pin 1  (3.3V)
GND   -----> Pin 6  (GND)
DATA  -----> Pin 11 (GPIO 17) ← MUDAR DE PIN 7
```

---

## 📝 SOLUÇÃO MANUAL (Se preferir)

```bash
# 1. Editar config
nano /opt/mysensoringo/internal/config/config.go

# 2. Mudar linha:
# De: GPIO: "4"
# Para: GPIO: "17"

# 3. Salvar (Ctrl+O, Enter, Ctrl+X)

# 4. Rebuild
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## ✅ VERIFICAR SE FUNCIONOU

```bash
sudo journalctl -u mysensoringo -f
```

**Se funcionar, deve ver:**
```
Metric saved: Temp=24.5°C, Humidity=65.0%
```

**Se ainda der erro:**
- Execute `bash scripts/fix-gpio.sh` novamente
- Escolha outro pino (18, 27 ou 22)

---

## 📚 Informação Completa

Veja `GPIO_SOLUTION.md` para:
- Explicação detalhada do problema
- Lista de pinos e status
- Alternativa de desabilitar 1-Wise

---

## 🎯 Resumo

**Agora mesmo:** Execute o script e mude o conector físico:
```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
# Mova fio DATA de Pin 7 para Pin 11
```
