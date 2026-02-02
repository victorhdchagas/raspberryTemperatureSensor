# MySensorIngo - Atualização Importante

## ⚠️ GPIO 4 e Problemas de Permissão

### O Problema
O GPIO 4 (BCM) no Raspberry Pi conflita com o sistema 1-Wire. Por padrão, o 1-Wire pode estar habilitado, bloqueando o uso do GPIO 4 para outros dispositivos.

### ✅ Solução Aplicada

**O serviço agora roda como `root`** para garantir acesso ao GPIO. Isso resolve a maioria dos problemas de permissão.

**Rebuild necessário:**
```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## 🔧 Configuração Avançada

### Mudar o Pino GPIO

Se ainda assim não funcionar com o pino 4, mude para outro pino:

```bash
# Opção 1: Variável de ambiente
export SENSOR_GPIO_PIN=17
bash scripts/build-local.sh

# Opção 2: Editar config.go
nano internal/config/config.go
# Mude GPIO: "4" para GPIO: "17"
```

### Pinos Recomendados

| BCM Number | Physical Pin | Uso |
|-------------|---------------|------|
| 4 | Pin 7 | DHT11 (padrão) - pode conflitar com 1-Wire |
| 17 | Pin 11 | ✅ Recomendado - sem conflito |
| 18 | Pin 12 | ✅ Alternativa |
| 27 | Pin 13 | ✅ Alternativa |
| 22 | Pin 15 | ✅ Alternativa |

---

## 📚 Documentação Completa

- **Instalação:** `DEPLOY_LOCAL.md`
- **Troubleshooting:** `TROUBLESHOOTING.md`
- **Deploy remoto:** `scripts/deploy.sh`

---

## 🚀 Deploy Rápido (Recomendado)

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

---

## 🔍 Verificação

Após o deploy, verifique se funcionou:

```bash
# Ver status do serviço
sudo systemctl status mysensoringo

# Ver logs
sudo journalctl -u mysensoringo -f

# Deve ver algo como:
# 2026/02/02 XX:XX:XX Metric saved: Temp=24.5°C, Humidity=65.0%
```

Se ainda houver erro "failed to export pin", consulte `TROUBLESHOOTING.md`.
