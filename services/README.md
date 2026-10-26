# Shared development services

Docker Compose stack for services shared across the apps you develop on this
machine. Runs inside colima — see the top-level README, "Docker (colima)".

## Stack

| Service | Image | Host port | Notes |
|---|---|---|---|
| redis | `redis:8.10` | 6379 | shared cache; persisted |
| mariadb | `mariadb:11.8` (LTS) | — (internal only) | Matomo's DB; Matomo supports MySQL/MariaDB only, never Postgres |
| matomo | `matomo:5.14.0-apache` | 8090 | self-contained PHP + Apache |

## Usage

```sh
cd ~/.config/services
cp .env.example .env    # once: set MATOMO_DB_PASSWORD
./services up           # everything
./services up matomo    # matomo + mariadb as a pair
./services stop matomo  # stops both
./services logs matomo
./services status
```

## matomo ↔ mariadb lifecycle tie

Compose `depends_on` only controls startup. The `services` helper in this
directory expands a `matomo` target to `matomo mariadb` on both `up` and
`stop`, so the database follows Matomo's lifecycle. Redis keeps its own
lifecycle — it is shared and independent of Matomo.

## Storage

All state is bind-mounted under `$XDG_STATE_HOME/services/`
(`~/.local/state/services/`), not named volumes. It stays visible on the
host and survives `colima delete` / VM re-provisioning (named volumes live
inside the VM disk).

```
~/.local/state/services/
├── mariadb/   # MariaDB data directory
├── matomo/    # /var/www/html (config, plugins, uploads)
└── redis/     # RDB snapshots
```

Redis persistence is on by default for uniformity. To make redis a pure
ephemeral cache, delete its `volumes:` line in `compose.yaml`.

First Matomo run: open http://127.0.0.1:8090 and complete the setup wizard
(the DB connection comes from `.env`; the wizard creates the admin account).

## Adding a service

Add a block to `compose.yaml`: pin the image tag, and persist state with a
bind mount under `$XDG_STATE_HOME/services/<name>`. Example (Postgres —
use this for your other apps; Matomo cannot share it):

```yaml
  postgres:
    image: postgres:17
    restart: unless-stopped
    ports:
      - "127.0.0.1:5432:5432"
    environment:
      POSTGRES_DB: ${POSTGRES_DB_NAME:-app}
      POSTGRES_USER: ${POSTGRES_USER:-app}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD:?set in .env}
    volumes:
      - ${XDG_STATE_HOME:?set XDG_STATE_HOME}/services/postgres:/var/lib/postgresql/data
```
