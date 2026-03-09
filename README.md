# deckhand

A lightweight, rsync-based backup tool for Docker environments.

Deckhand discovers containers via Docker labels, optionally stops them before syncing, and rsyncs their data to a remote destination (e.g. a NAS over SSH). No tarballs, no proprietary formats — just files in the right place.

## How it works

1. Discovers running containers with `deckhand.enable=true`
2. Stops containers marked with `deckhand.stop=true` (in priority order)
3. Runs any `deckhand.pre-exec` hooks
4. rsyncs configured paths to the destination
5. Restarts stopped containers (always — even on error)
6. Optionally fires a webhook with the result

## Container labels

| Label | Required | Description |
|---|---|---|
| `deckhand.enable` | yes | Set to `true` to include this container |
| `deckhand.paths` | yes | Comma-separated host paths to sync |
| `deckhand.stop` | no | Stop container during sync (default: `false`) |
| `deckhand.pre-exec` | no | Command to run inside the container before sync |
| `deckhand.priority` | no | Stop/start order — lower number = stopped first (default: `0`) |

## Configuration

| Environment variable | Default | Description |
|---|---|---|
| `BACKUP_DESTINATION` | *(required)* | rsync destination, e.g. `user@nas:/mnt/backups/docker` |
| `BACKUP_CRON` | `0 2 * * *` | Cron schedule |
| `BACKUP_SSH_KEY` | `/keys/id_rsa` | Path to SSH private key |
| `BACKUP_RSYNC_ARGS` | | Extra args passed to rsync |
| `BACKUP_NOTIFY_URL` | | Webhook URL called on completion/failure |
| `BACKUP_DRY_RUN` | `false` | Run rsync with `--dry-run` |
| `BACKUP_RUN_ONCE` | `false` | Run once and exit (useful for testing) |

## Usage

```yaml
services:
  deckhand:
    image: ghcr.io/olvrvrmr/deckhand:latest
    container_name: deckhand
    restart: unless-stopped
    environment:
      - BACKUP_DESTINATION=user@nas:/mnt/backups/docker
      - BACKUP_SSH_KEY=/keys/id_rsa
      - BACKUP_CRON=0 2 * * *
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - /path/to/ssh/key:/keys/id_rsa:ro

  myapp:
    image: myapp:latest
    labels:
      - "deckhand.enable=true"
      - "deckhand.stop=true"
      - "deckhand.paths=/opt/appdata/myapp"
      - "deckhand.priority=10"
```

## Recovery

Restoring is as simple as rsyncing back from the destination and starting your containers:

```bash
rsync -av user@nas:/mnt/backups/docker/myapp/ /opt/appdata/myapp/
docker compose up -d
```

## Socket proxy

If you use a Docker socket proxy (recommended), ensure the following permissions are enabled:

```
CONTAINERS=1
INFO=1
```
