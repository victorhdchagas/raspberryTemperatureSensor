# Solução de Problemas GPIO - Raspberry Pi

## Erro: "failed to export pin 4"

Este erro ocorre quando há problema de permissão ou conflito com o sistema 1-Wire do Raspberry Pi.

---

## Solução 1: Rodar como Root ✅ (Implementado)

O serviço agora roda como `root` para garantir acesso ao GPIO. Rebuild necessário:

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## Solução 2: Desabilitar 1-Wire (Se usar GPIO 4)

O GPIO 4 (BCM) é usado por padrão pelo sistema 1-Wire. Para usar o DHT11 no pino 4:

```bash
# Editar config.txt
sudo nano /boot/config.txt

# Adicionar ou descomentar:
dtoverlay=w1-gpio-pullup,gpiopin=4

# OU desabilitar completamente (se não usar 1-Wire):
# dtoverlay=w1-gpio-pullup,gpiopin=4

# Salvar e reiniciar
sudo reboot
```

**Recomendação:** Se o teste Python funciona com `board.D4`, então o 1-Wire está desabilitado. O problema é apenas permissão → Solução 1 deve resolver.

---

## Solução 3: Mudar para Outro Pino

Se não conseguir usar o pino 4, mude para outro pino (ex: GPIO 17):

1. Editar `internal/config/config.go`:
```go
Sensor: SensorConfig{
    GPIO:     "17",  // Mudar de "4" para "17"
    Interval: 50 * time.Minute,
},
```

2. Rebuild:
```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

**Pinos recomendados (BCM numbering):**
- GPIO 17 (physical pin 11)
- GPIO 18 (physical pin 12)
- GPIO 27 (physical pin 13)
- GPIO 22 (physical pin 15)

---

## Verificação Rápida

Testar acesso GPIO manualmente:

```bash
# Verificar permissões
ls -l /sys/class/gpio/

# Se necessário, adicionar usuário ao grupo gpio (se não usar root)
sudo usermod -a -G gpio pi
```

---

## Pinos GPIO vs Physical

| BCM Number | Physical Pin | Descrição |
|-------------|---------------|------------|
| 4 | Pin 7 | DHT11 (padrão) - conflita com 1-Wire |
| 17 | Pin 11 | ✅ Alternativa segura |
| 18 | Pin 12 | ✅ Alternativa segura |
| 27 | Pin 13 | ✅ Alternativa segura |
| 22 | Pin 15 | ✅ Alternativa segura |

---

## Teste Rápido com Diferentes Pinós

Para testar rapidamente outros pinós:

1. Conecte DHT11 Data pin no GPIO 17 (physical pin 11)
2. Edite config (mude GPIO de "4" para "17")
3. Rebuild:
```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## Diagnóstico Completo

Se ainda assim não funcionar, verifique:

```bash
# Verificar se 1-Wire está habilitado
ls /sys/bus/w1/devices/

# Verificar se há dispositivos 1-Wire (indicando conflito)
cat /sys/bus/w1/w1_bus_master1/available_w1_bus_masters

# Testar export manual do GPIO
echo 4 | sudo tee /sys/class/gpio/export  # Deve funcionar
ls /sys/class/gpio/gpio4/  # Verificar se criou
echo out | sudo tee /sys/class/gpio/gpio4/direction
echo 4 | sudo tee /sys/class/gpio/unexport
```

---

## Resumo

1. **Rodar como root** ✅ (já configurado)
2. **Se não funcionar:** Desabilitar 1-Wire ou mudar pino
3. **Recomendação:** GPIO 17 se o 4 continuar com problemas

Execute o rebuild após qualquer mudança de configuração!
