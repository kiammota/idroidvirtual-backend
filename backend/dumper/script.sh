#!/bin/bash

#############################################
# Backup automático PostgreSQL - Idroid
#############################################

# === CONFIGURAÇÕES ===
CONTAINER="idroid-db"                 # nome do container do banco
DB_USER="vps"
DB_PASSWORD="vps123"
DB_NAME="idroid_db"

BACKUP_DIR="/home/vps/backups/postgres"
RETENTION_DAYS=14
LOG_FILE="/var/log/backup-idroid.log"

# === GERA O NOME DO ARQUIVO ===
case $(date +%u) in
  1) DIA="SEGUNDA-FEIRA" ;;
  2) DIA="TERCA-FEIRA" ;;
  3) DIA="QUARTA-FEIRA" ;;
  4) DIA="QUINTA-FEIRA" ;;
  5) DIA="SEXTA-FEIRA" ;;
  6) DIA="SABADO" ;;
  7) DIA="DOMINGO" ;;
esac

case $(date +%m) in
  01) MES="JANEIRO" ;;
  02) MES="FEVEREIRO" ;;
  03) MES="MARCO" ;;
  04) MES="ABRIL" ;;
  05) MES="MAIO" ;;
  06) MES="JUNHO" ;;
  07) MES="JULHO" ;;
  08) MES="AGOSTO" ;;
  09) MES="SETEMBRO" ;;
  10) MES="OUTUBRO" ;;
  11) MES="NOVEMBRO" ;;
  12) MES="DEZEMBRO" ;;
esac

DIA_NUM=$(date +%d)
HORA=$(date +%Hh%M)

FILENAME="DUMP-${DIA}-${DIA_NUM}-${MES}-${HORA}.sql.gz"

# === PREPARAÇÃO ===
mkdir -p "$BACKUP_DIR"

# === EXECUTA O BACKUP VIA DOCKER ===
docker exec -e PGPASSWORD="$DB_PASSWORD" "$CONTAINER" \
  pg_dump -U "$DB_USER" -d "$DB_NAME" \
  --no-owner \
  --no-acl \
  --clean \
  --if-exists | gzip > "$BACKUP_DIR/$FILENAME"

# === VERIFICAÇÃO ===
if [ $? -eq 0 ]; then
  echo "$(date '+%Y-%m-%d %H:%M:%S') - Backup OK → $FILENAME" >> "$LOG_FILE"
  echo "Backup criado com sucesso: $BACKUP_DIR/$FILENAME"
else
  echo "$(date '+%Y-%m-%d %H:%M:%S') - ERRO no backup!" >> "$LOG_FILE"
  echo "Erro ao criar o backup"
  exit 1
fi

# === LIMPEZA DE BACKUPS ANTIGOS ===
find "$BACKUP_DIR" -name "DUMP-*.sql.gz" -type f -mtime +$RETENTION_DAYS -delete
