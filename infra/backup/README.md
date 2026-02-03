# Резервное копирование Вокзал.ТЕХ

Скрипты и рекомендации по резервному копированию в соответствии с [docs/initial/20.md](../../docs/initial/20.md). Цель — сохранность данных, быстрое восстановление и соответствие 152-ФЗ (хранение логов и данных 1 год).

## Что резервируется

| Компонент    | Частота      | Метод              | Хранение   |
|-------------|--------------|--------------------|------------|
| PostgreSQL  | Каждые 4 ч   | pg_dump (см. ниже) | MinIO/S3   |
| Redis       | Каждые 6 ч   | RDB + AOF          | MinIO      |
| Документы   | По требованию| rsync / mc         | MinIO/S3   |
| Конфигурации| При изменении| Git                | GitHub     |

## Скрипт backup-postgres.sh

Создаёт дамп БД в формате custom (`pg_dump -Fc`), при необходимости шифрует его GPG и удаляет старые дампы по сроку хранения.

### Требования

- Доступ к PostgreSQL (клиент `pg_dump` в PATH).
- Переменные окружения для подключения к БД (или ~/.pgpass): `PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`.

### Использование

```bash
cd infra/backup
chmod +x backup-postgres.sh

# Базовый запуск (дамп в ./backups)
./backup-postgres.sh

# С указанием каталога и срока хранения
BACKUP_DIR=/var/backups/vokzal RETENTION_DAYS=365 ./backup-postgres.sh

# С шифрованием GPG
GPG_RECIPIENT=backup@vokzal.tech ./backup-postgres.sh
```

### Переменные окружения

| Переменная     | Описание                          | По умолчанию |
|----------------|-----------------------------------|--------------|
| BACKUP_DIR     | Каталог для сохранения дампов     | ./backups    |
| RETENTION_DAYS | Хранить дампы (дней)              | 365          |
| PGDATABASE     | Имя базы данных                   | vokzal       |
| GPG_RECIPIENT  | E-mail ключа GPG для шифрования   | —            |

## Запуск по расписанию (cron)

Пример — полный дамп каждые 4 часа и очистка старых:

```cron
0 */4 * * * BACKUP_DIR=/var/backups/vokzal PGPASSWORD=xxx /path/to/infra/backup/backup-postgres.sh
```

Ежедневный дамп в 02:00 с шифрованием:

```cron
0 2 * * * BACKUP_DIR=/var/backups/vokzal GPG_RECIPIENT=backup@vokzal.tech /path/to/infra/backup/backup-postgres.sh
```

## Выгрузка в MinIO/S3

После создания дампа его можно загрузить в MinIO или S3-совместимое хранилище (см. doc 20):

- Установите [MinIO Client (mc)](https://min.io/docs/minio/linux/reference/minio-mc.html) и настройте alias.
- Пример загрузки последнего дампа:

```bash
mc cp "${BACKUP_DIR}"/vokzal_*.dump myminio/vokzal-backups/
```

Либо используйте `aws s3 cp` / скрипт-обёртку с переменными `S3_BUCKET`, `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`.

## Восстановление PostgreSQL

1. Скачайте нужный дамп из хранилища (например, `latest.dump` или по дате).
2. При необходимости расшифруйте:  
   `gpg --decrypt latest.dump.gpg > latest.dump`
3. Восстановите базу (существующую БД перезапишет):

```bash
pg_restore -U vokzal -d vokzal_db --clean --if-exists latest.dump
```

Или для создания новой БД:

```bash
createdb -U vokzal vokzal_restored
pg_restore -U vokzal -d vokzal_restored latest.dump
```

## Redis

Резервное копирование Redis выполняется средствами самого Redis (RDB, AOF). Настройте в `redis.conf`:

- `save` или `appendonly yes` и периодическое копирование файлов RDB/AOF в MinIO или на другой сервер по расписанию (cron + mc/rsync).

## Тестирование восстановления

Рекомендуется раз в квартал выполнять полное восстановление в тестовую среду и проверять целостность данных и работу продажи/возврата (в т.ч. соответствие 54-ФЗ).

## Контакты

- Резервное хранилище: backup.vokzal.tech (по настройкам окружения)
- Документация: [docs/initial/20.md](../../docs/initial/20.md)

---

© 2026 Вокзал.ТЕХ
