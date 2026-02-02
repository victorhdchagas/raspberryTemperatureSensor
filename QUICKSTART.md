# Deploy Rápido - MySensorIngo

## Script Corrigido ✅

O problema "Text file busy" foi resolvido. O script agora:
1. **Para o serviço** antes de copiar o novo binário
2. Copia os arquivos
3. **Reinicia o serviço** automaticamente

## Como Executar

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

## O que acontece:

1. `git pull` - Atualiza código
2. Build - Gera `mysensoringo`
3. `systemctl stop` - Para serviço (evita erro "Text file busy")
4. `cp` - Copia novo binário
5. `systemctl start` - Inicia serviço com leitura imediata!

## Verificar se funcionou:

```bash
# Ver status
sudo systemctl status mysensoringo

# Ver logs
sudo journalctl -u mysensoringo -f

# Deve ver:
# Metric saved: Temp=XX.X°C, Humidity=XX.X%
```

## Acesso Dashboard

```
http://<ip-do-raspberry>:8080
```

## Atualizações Futuras

Sempre que quiser atualizar, basta executar o mesmo script:

```bash
cd /opt/mysensoringo
bash scripts/build-local.sh
```

O script cuida de tudo! 🚀
