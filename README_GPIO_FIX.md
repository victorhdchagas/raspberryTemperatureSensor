# MySensorIngo - Problema GPIO 4 Resolvido

## 🚨 Erro: "failed to export pin 4"

**O problema não é permissão, é CONFLITO HARDWARE.**

O GPIO 4 (BCM) é usado automaticamente pelo sistema 1-Wise do Raspberry Pi.

---

## ✅ SOLUÇÃO APLICADA

1. **Script de teste automático** - `scripts/fix-gpio.sh`
   - Testa múltiplos pinos
   - Mostra quais funcionam
   - Muda configuração automaticamente

2. **Documentação completa** - `GPIO_SOLUTION.md`
   - Explicação detalhada
   - Diagrama de pinos
   - Alternativas

3. **Instruções rápidas** - `SOLUTION_NOW.md`
   - O que fazer AGORA
   - Diagrama de conexão física

---

## 🚀 Para Resolver Agora

### No Raspberry Pi:

```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
```

### Conectar DHT11 no GPIO 17:

```
VCC → Pin 1 (3.3V)
GND → Pin 6 (GND)
DATA → Pin 11 (GPIO 17) ← MUDAR DE PIN 7
```

### Verificar funcionamento:

```bash
sudo journalctl -u mysensoringo -f
# Deve ver: Metric saved: Temp=XX.X°C, Humidity=XX.X%
```

---

## 📌 Pinos Recomendados

| GPIO | Physical | Funciona? | Nota |
|-------|----------|------------|-------|
| 4 | Pin 7 | ❌ Não - conflita 1-Wise |
| 17 | Pin 11 | ✅ Sim - **Recomendado** |
| 18 | Pin 12 | ✅ Sim - mas pode conflitar áudio |
| 27 | Pin 13 | ✅ Sim |
| 22 | Pin 15 | ✅ Sim |

---

## 📚 Documentação

- `SOLUTION_NOW.md` - O que fazer AGORA
- `GPIO_SOLUTION.md` - Solução completa e detalhada
- `scripts/fix-gpio.sh` - Script de teste automático
- `DEPLOY_LOCAL.md` - Deploy local
- `TROUBLESHOOTING.md` - Outros problemas

---

## ✅ Status

**Correção aplicada:** Script de teste + documentação completa
**Próximo passo:** Usar script `fix-gpio.sh` no Raspberry Pi
