# PostgreSQL init scripts

Scripts in this directory run once when the PostgreSQL container is first created (see `docker-compose.yml` → postgres → volumes → `./init-scripts:/docker-entrypoint-initdb.d`).

## 01-init.sh

Creates:

- Extension `uuid-ossp`
- One schema per microservice: `auth`, `schedule`, `ticket`, `fiscal`, `payment`, `board`, `geo`, `notify`, `audit`, `document`

## Single-admin user (intentional for Docker)

All microservices are configured to connect as the same PostgreSQL user (`admin` in `docker-compose.yml`). The init script does **not** create separate DB users (e.g. `vokzal_auth`, `vokzal_ticket`) or grant schema-specific privileges.

**Why:** Simplifies local development and CI: one user, one password, no per-service user/password management in configs.

**Production:** For production you must use the principle of least privilege:

- Create one PostgreSQL user per service (e.g. `vokzal_auth`, `vokzal_schedule`).
- Grant each user only `USAGE`, `CREATE` (and as needed `SELECT`/`INSERT`/…) on **its own** schema.
- Store credentials in secrets (e.g. Kubernetes Secrets) and configure each service with its dedicated DB user.

This repo does not include a multi-user init script; add one (or use a migration/operator) when deploying to production.
