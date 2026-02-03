#!/usr/bin/env bash
# Вокзал.ТЕХ — резервное копирование PostgreSQL (см. docs/initial/20.md)
# Использование: BACKUP_DIR=/path/to/backups ./backup-postgres.sh
# Переменные окружения (опционально):
#   PGHOST, PGPORT, PGUSER, PGPASSWORD, PGDATABASE — подключение к БД
#   BACKUP_DIR — каталог для дампов (по умолчанию ./backups)
#   RETENTION_DAYS — хранить дампы N дней (по умолчанию 365)
#   GPG_RECIPIENT — ключ для шифрования (если задан, дамп шифруется gpg)

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKUP_DIR="${BACKUP_DIR:-${SCRIPT_DIR}/backups}"
RETENTION_DAYS="${RETENTION_DAYS:-365}"
PGDATABASE="${PGDATABASE:-vokzal}"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"
DUMP_NAME="vokzal_${TIMESTAMP}.dump"
DUMP_PATH="${BACKUP_DIR}/${DUMP_NAME}"

mkdir -p "${BACKUP_DIR}"

echo "[$(date -Iseconds)] Starting PostgreSQL backup: ${PGDATABASE} -> ${DUMP_PATH}"
pg_dump -Fc -f "${DUMP_PATH}" "${PGDATABASE}"

if [[ -n "${GPG_RECIPIENT:-}" ]]; then
  gpg --encrypt --recipient "${GPG_RECIPIENT}" --trust-model always -o "${DUMP_PATH}.gpg" "${DUMP_PATH}"
  rm -f "${DUMP_PATH}"
  echo "[$(date -Iseconds)] Encrypted backup: ${DUMP_PATH}.gpg"
fi

echo "[$(date -Iseconds)] Removing backups older than ${RETENTION_DAYS} days"
find "${BACKUP_DIR}" -name "vokzal_*.dump*" -mtime "+${RETENTION_DAYS}" -delete 2>/dev/null || true

echo "[$(date -Iseconds)] Backup finished: ${DUMP_PATH}"
