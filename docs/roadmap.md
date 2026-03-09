# Deckhand Roadmap

Deckhand focuses on **simple, transparent Docker appdata backups using
rsync**.

The roadmap prioritizes:

-   reliability
-   predictable restore workflows
-   minimal configuration
-   homelab-friendly operation

Deckhand is intentionally designed to stay **small, predictable, and
easy to restore from**.

------------------------------------------------------------------------

# Phase 1 - Stabilization (v0.x → v1.0)

Goal: make Deckhand **rock solid and production ready**.

## Backup safety

-   Improve container stop/restart reliability
-   Ensure containers restart even if rsync fails
-   Handle interrupted backups gracefully
-   Improve error reporting

------------------------------------------------------------------------

## Logging improvements

Add better logging and debugging tools.

Examples:

-   structured logs
-   verbose/debug mode
-   clearer rsync output

Example environment variable:

    BACKUP_LOG_LEVEL=debug

------------------------------------------------------------------------

## Label validation

Detect common configuration mistakes early:

-   missing `deckhand.path`
-   invalid paths
-   unknown labels

Deckhand should fail early with **clear error messages**.

------------------------------------------------------------------------

## Restore documentation

Add official restore workflows:

-   restoring a single container
-   restoring a full Docker stack
-   restoring from NAS backups

------------------------------------------------------------------------

## CI Improvements

Introduce GitHub Actions for:

-   linting
-   ✅ container build
-   ✅ automated releases (triggered by git tag push)
-   ✅ version tagging (manual tag → automated image tags + GitHub Release)

------------------------------------------------------------------------

## Documentation structure

Introduce structured documentation:

    docs/
     ├─ labels.md
     ├─ architecture.md
     └─ roadmap.md

------------------------------------------------------------------------

# Phase 2 - Power Features (v1.x)

These features make Deckhand more powerful while keeping the project
simple.

------------------------------------------------------------------------

## Incremental snapshot backups

Support rsync snapshot-style backups using `--link-dest`.

Example structure:

    backups/
     ├─ 2026-03-01/
     ├─ 2026-03-02/
     ├─ 2026-03-03/
     └─ latest/

Advantages:

-   efficient storage
-   easy rollback
-   still plain files

------------------------------------------------------------------------

## Container pause mode

Instead of stopping containers:

    docker pause
    rsync
    docker unpause

This reduces downtime for some services.

------------------------------------------------------------------------

## Backup metrics ✅

Expose Prometheus metrics such as:

-   ✅ `deckhand_backups_total`
-   ✅ `deckhand_backup_duration_seconds`
-   ✅ `deckhand_backup_failures_total`
-   ✅ `deckhand_backup_running`
-   ✅ `deckhand_last_backup_status`
-   ✅ `deckhand_last_backup_timestamp`
-   ✅ `deckhand_bytes_transferred_total`
-   ✅ `deckhand_containers_discovered`

See [metrics.md](metrics.md) for full reference and included Grafana dashboard.

------------------------------------------------------------------------

# Phase 3 - Long-term Ideas

These are potential improvements but not core goals.

------------------------------------------------------------------------

## Optional Web UI

A minimal dashboard showing:

-   containers discovered
-   last backup status
-   backup history

> Considered once the core tool is stable. Lower priority than reliability and observability features.

------------------------------------------------------------------------

## Backup policies

Labels could define backup policies:

    deckhand.schedule=nightly
    deckhand.retention=7d

------------------------------------------------------------------------

## Multi-destination backups

Allow sending backups to multiple destinations:

    NAS
    +
    offsite backup

------------------------------------------------------------------------

## Plugin system

Allow optional backup extensions:

-   database plugins
-   cloud storage targets
-   compression options

------------------------------------------------------------------------

# Non‑Goals

Deckhand intentionally does **not** aim to become a full backup
platform.

It will **not replace** tools like:

-   Borg
-   Restic
-   Duplicati

Deckhand focuses exclusively on:

-   Docker appdata replication
-   simple rsync-based backups
-   easy restores

### Single path per container is intentional

Deckhand deliberately supports only one `deckhand.path` per container. Multiple paths per container introduce destination conflicts, inconsistent `--delete` behaviour, and complicate restore workflows. For multi-container apps (e.g. app + database), the recommended pattern is to enable backup on one container only and use `deckhand.stop=true` without a path on dependents.
