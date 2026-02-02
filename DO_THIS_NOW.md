# 🎯 AÇÃO AGORA - Resolver GPIO 4

## 🚨 O Problema

Erro: **"failed to export pin 4"**

**Causa:** GPIO 4 conflita com sistema 1-Wise do Raspberry Pi (conflito de hardware).
**Solução:** Mudar para GPIO 17.

---

## ✅ O QUE FAZER AGORA (3 passos)

### Passo 1: No Raspberry Pi - Executar build com GPIO 17

```bash
cd /opt/mysensoringo
bash scripts/build-gpio17.sh
```

### Passo 2: Mover o conector físico do DHT11

**Mova o fio DATA:**
- ❌ De: **Pin 7** (GPIO 4)
- ✅ Para: **Pin 11** (GPIO 17)

```
DHT11       Raspberry Pi
---------   ---------------
VCC   -----> Pin 1  (3.3V)
GND   -----> Pin 6  (GND)
DATA  -----> Pin 11 (GPIO 17) ← MUDAR!
```

### Passo 3: Verificar se funcionou

```bash
sudo journalctl -u mysensoringo -f
```

**Se funcionou, deve ver:**
```
Metric saved: Temp=XX.X°C, Humidity=XX.X%
```

---

## 🔧 Se ainda não funcionar

Execute o script de teste de pinos:

```bash
cd /opt/mysensoringo
bash scripts/fix-gpio.sh
```

O script vai testar GPIO 4, 17, 18, 27, 22 e mudar automaticamente para um que funciona.

---

## 📚 Informação Completa

- `README.md` - Documentação principal atualizada
- `SOLUTION_NOW.md` - O que fazer AGORA
- `GPIO_SOLUTION.md` - Explicação completa
- `TROUBLESHOOTING.md` - Outros problemas
- `scripts/build-gpio17.sh` - Build com GPIO 17
- `scripts/fix-gpio.sh` - Testa pinos automaticamente

---

## ✅ Resumo

**Agora mesmo:**
1. Execute `bash scripts/build-gpio17.sh` no Raspberry Pi
2. Mova o fio DATA do DHT11 de Pin 7 para Pin 11
3. Verifique os logs: `sudo journalctl -u mysensoringo -f`

Pronto! 🚀
