# deckhand

![GitHub
release](https://img.shields.io/github/v/release/olvrvrmr/deckhand)
![License](https://img.shields.io/github/license/olvrvrmr/deckhand)

**Label-driven Docker appdata backups with rsync.**

Docker containers are easy to redeploy, but their persistent data is
not.

Deckhand discovers Docker containers via labels, optionally stops them
for consistency, and rsyncs their persistent data to a remote
destination such as a NAS over SSH.

No proprietary backup format. No archives. No complex restore process.\
Just your files, synced safely and ready to restore.

------------------------------------------------------------------------

## Why Deckhand?

Deckhand is built for people who want a simple and predictable way to
back up Docker appdata:

-   **Docker-native**: opt containers in with labels
-   **rsync-based**: backups stay readable and portable
-   **restore-friendly**: sync files back and start your containers
-   **consistency-aware**: optionally stop containers before backup
-   **homelab-friendly**: ideal for appdata-to-NAS workflows

------------------------------------------------------------------------

## Architecture

Docker Host\
│\
│ container labels\
▼\
Deckhand container\
│\
│ rsync\
▼\
NAS / remote storage

------------------------------------------------------------------------

## Quick start

Add Deckhand to your Docker host:

``` yaml
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
      - /mnt/appdata:/mnt/appdata
      - /path/to/ssh/key:/keys/id_rsa:ro
```

Then label any container you want to back up:

``` yaml
services:
  myapp:
    image: myapp:latest
    labels:
      - "deckhand.enable=true"
      - "deckhand.stop=true"
      - "deckhand.path=/mnt/appdata/myapp"
      - "deckhand.exclude=logs,*.tmp"
```

That's it. Deckhand will discover the container and back up the
configured path on schedule.

------------------------------------------------------------------------

## How it works

1.  Discovers containers with `deckhand.enable=true`
2.  Stops containers marked with `deckhand.stop=true`
3.  Runs any `deckhand.pre-exec` hooks
4.  rsyncs configured paths to the backup destination
5.  Restarts stopped containers, even if an error occurs
6.  Optionally sends a webhook notification

------------------------------------------------------------------------

## Label reference

Deckhand is controlled entirely via container labels.

  Label                 Description
  --------------------- ----------------------------------------------
  `deckhand.enable`     Enables backup for the container
  `deckhand.path`       Path inside the host filesystem to back up
  `deckhand.stop`       Stop container before backup for consistency
  `deckhand.exclude`    Comma-separated list of patterns to exclude
  `deckhand.pre-exec`   Command to run before backup
  `deckhand.priority`   Backup order priority

Example:

``` yaml
labels:
  - "deckhand.enable=true"
  - "deckhand.path=/mnt/appdata/myapp"
  - "deckhand.stop=true"
  - "deckhand.exclude=cache,tmp"
```

------------------------------------------------------------------------

## Typical use cases

Deckhand is ideal for:

-   Homelab Docker servers backing up appdata to a NAS
-   Self-hosted environments where restore simplicity is important
-   Systems where traditional backup tools are too heavy
-   Environments where Docker containers should be backed up
    automatically via labels

------------------------------------------------------------------------

## Restore

Because Deckhand stores backups as plain files, restoring is just an
rsync away.

``` bash
rsync -av user@nas:/mnt/backups/docker/myapp/ /mnt/appdata/myapp/
docker compose up -d
```

------------------------------------------------------------------------

## License

MIT
