# Deckhand Label Specification

Deckhand uses Docker container labels to determine which containers
should be backed up and how their data should be handled.

All labels are namespaced under:

    deckhand.*

Only containers with `deckhand.enable=true` are considered by Deckhand.

------------------------------------------------------------------------

# Required Label

## `deckhand.enable`

Enables Deckhand backup for the container.

**Type**

    boolean

**Example**

``` yaml
labels:
  - "deckhand.enable=true"
```

If this label is not present, the container is ignored.

------------------------------------------------------------------------

# Core Backup Labels

## `deckhand.path`

Defines the **host path** that should be backed up.

This should match the host path where the container stores its
persistent data.

**Type**

    string

**Example**

``` yaml
labels:
  - "deckhand.enable=true"
  - "deckhand.path=/mnt/appdata/sonarr"
```

Multiple paths may be supported in the future.

------------------------------------------------------------------------

## `deckhand.exclude`

Comma-separated list of patterns that should be excluded from the
backup.

These are passed directly to `rsync --exclude`.

**Type**

    string (comma-separated)

**Example**

``` yaml
labels:
  - "deckhand.exclude=logs,*.tmp,cache"
```

Equivalent rsync flags:

    --exclude logs
    --exclude *.tmp
    --exclude cache

------------------------------------------------------------------------

# Consistency Labels

## `deckhand.stop`

Stops the container before backup and restarts it afterward.

Useful for services that maintain open databases (SQLite, etc.).

**Type**

    boolean

**Default**

    false

**Example**

``` yaml
labels:
  - "deckhand.stop=true"
```

Backup flow:

    docker stop
    rsync backup
    docker start

------------------------------------------------------------------------

## `deckhand.pre-exec`

Command executed **before the backup begins**.

This is useful for:

-   database dumps
-   flushing caches
-   application-specific backup routines

**Type**

    string

**Example**

``` yaml
labels:
  - "deckhand.pre-exec=pg_dumpall -U postgres > /backup/db.sql"
```

Execution occurs inside the container using `docker exec`.

------------------------------------------------------------------------

# Backup Ordering

## `deckhand.priority`

Defines the **backup order**.

Containers with lower numbers are backed up first.

Useful when dependencies exist between services.

**Type**

    integer

**Default**

    100

**Example**

``` yaml
labels:
  - "deckhand.priority=10"
```

------------------------------------------------------------------------

# Optional Labels (Future Expansion)

These labels are not yet implemented but reserved for future
functionality.

------------------------------------------------------------------------

## `deckhand.post-exec`

Command executed **after the backup completes**.

Possible use cases:

-   database cleanup
-   application restart hooks
-   cache rebuilds

Example:

``` yaml
labels:
  - "deckhand.post-exec=echo backup complete"
```

------------------------------------------------------------------------

## `deckhand.pause`

Pauses the container instead of stopping it.

Potential implementation:

    docker pause
    rsync
    docker unpause

This could reduce downtime for certain applications.

------------------------------------------------------------------------

## `deckhand.snapshot`

Enables timestamped snapshot backups.

Instead of overwriting backups, Deckhand would create:

    backup/
     ├─ 2026-03-07/
     ├─ 2026-03-08/
     └─ 2026-03-09/

Using `rsync --link-dest`.

------------------------------------------------------------------------

# Complete Example

``` yaml
services:
  sonarr:
    image: linuxserver/sonarr
    volumes:
      - /mnt/appdata/sonarr:/config
    labels:
      - "deckhand.enable=true"
      - "deckhand.path=/mnt/appdata/sonarr"
      - "deckhand.stop=true"
      - "deckhand.exclude=logs,*.tmp"
      - "deckhand.priority=20"
```

------------------------------------------------------------------------

# Backup Lifecycle

When Deckhand runs, the following sequence occurs:

    1 Discover containers with deckhand.enable=true
    2 Sort containers by deckhand.priority
    3 Execute deckhand.pre-exec (if present)
    4 Stop container if deckhand.stop=true
    5 Run rsync backup
    6 Restart container if it was stopped
    7 Send webhook notification (optional)

------------------------------------------------------------------------

# Design Principles

Deckhand labels follow these principles:

### Explicit opt-in

Containers are only backed up if:

    deckhand.enable=true

------------------------------------------------------------------------

### Transparent backups

Backups are stored as normal files, allowing:

    rsync restore
    manual inspection
    filesystem snapshots

------------------------------------------------------------------------

### Minimal configuration

Most containers only need:

    deckhand.enable
    deckhand.path
