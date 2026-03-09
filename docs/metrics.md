# Deckhand Metrics Specification

Deckhand exposes Prometheus-compatible metrics to allow monitoring of
backup activity, failures, performance, and data transfer.

Metrics are exposed through an HTTP endpoint:

    /metrics

Typical endpoint:

    http://deckhand:2112/metrics

Port **2112** is commonly used for Prometheus exporters.

------------------------------------------------------------------------

# Metric Design Principles

Deckhand metrics follow these principles:

-   **Low cardinality**
-   **Clear container-level visibility**
-   **Minimal but useful set of metrics**
-   **Prometheus naming conventions**
-   **Operational insight for homelabs**

Metrics are labeled by container name where appropriate.

------------------------------------------------------------------------

# Metrics

## deckhand_backups_total

Total number of backup attempts.

**Type**

    Counter

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container
  status      success or failure

**Example**

    deckhand_backups_total{container="sonarr",status="success"} 12
    deckhand_backups_total{container="radarr",status="failure"} 2

------------------------------------------------------------------------

## deckhand_backup_failures_total

Total number of failed backups.

**Type**

    Counter

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Example**

    deckhand_backup_failures_total{container="nextcloud"} 3

------------------------------------------------------------------------

## deckhand_backup_duration_seconds

Duration of each backup execution.

**Type**

    Histogram

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Example**

    deckhand_backup_duration_seconds{container="plex"} 8.24

This metric helps identify:

-   slow NAS performance
-   large datasets
-   potential network bottlenecks

------------------------------------------------------------------------

## deckhand_last_backup_timestamp

Unix timestamp of the last successful backup.

**Type**

    Gauge

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Example**

    deckhand_last_backup_timestamp{container="sonarr"} 1710192000

This metric enables alerts when backups become stale.

Example Prometheus query:

    time() - deckhand_last_backup_timestamp

------------------------------------------------------------------------

## deckhand_bytes_transferred_total

Total amount of data transferred by backups.

**Type**

    Counter

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Example**

    deckhand_bytes_transferred_total{container="plex"} 1298391823

This metric can be used to observe:

-   growth of application data
-   network usage during backups

------------------------------------------------------------------------

## deckhand_backup_running

Indicates whether a backup is currently in progress for a container.

**Type**

    Gauge

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Values**

  Value   Meaning
  ------- -----------------------
  1       Backup in progress
  0       Idle

**Example**

    deckhand_backup_running{container="sonarr"} 1

------------------------------------------------------------------------

## deckhand_last_backup_status

Result of the last completed backup attempt for a container.

**Type**

    Gauge

**Labels**

  Label       Description
  ----------- ------------------------------
  container   Name of the Docker container

**Values**

  Value   Meaning
  ------- -----------------------
  1       Last attempt succeeded
  0       Last attempt failed

**Example**

    deckhand_last_backup_status{container="sonarr"} 1
    deckhand_last_backup_status{container="radarr"} 0

This metric always reflects the most recent run, regardless of outcome.
Combined with `deckhand_last_backup_timestamp`, it provides a complete
picture: the timestamp shows freshness of the last valid restore point,
while this metric shows whether the latest attempt succeeded.

------------------------------------------------------------------------

## deckhand_containers_discovered

Number of containers with `deckhand.enable=true` and a configured backup
path, as discovered in the last backup run. Stop-only containers (no
path defined) are excluded.

**Type**

    Gauge

**Example**

    deckhand_containers_discovered 5

Use this as a sanity check against `deckhand_backups_total{status="success"}`:
after a clean run both values should match.

------------------------------------------------------------------------

# Example Prometheus Configuration

Example scrape configuration:

``` yaml
scrape_configs:
  - job_name: deckhand
    static_configs:
      - targets: ["deckhand:2112"]
```

------------------------------------------------------------------------

# Example Prometheus Alerts

## Backup failure alert

``` yaml
- alert: DeckhandBackupFailure
  expr: increase(deckhand_backup_failures_total[1h]) > 0
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: Deckhand backup failure detected
```

------------------------------------------------------------------------

## Stale backup alert

``` yaml
- alert: DeckhandBackupStale
  expr: time() - deckhand_last_backup_timestamp > 86400
  for: 10m
  labels:
    severity: critical
  annotations:
    summary: Deckhand backup has not run in 24 hours
```

------------------------------------------------------------------------

## Last backup failed alert

``` yaml
- alert: DeckhandLastBackupFailed
  expr: deckhand_last_backup_status == 0
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Last backup failed for {{ $labels.container }}"
```

------------------------------------------------------------------------

# Example Grafana Panels

Typical Grafana dashboard panels:

-   Backup count per container
-   Backup duration
-   Backup failures
-   Last backup age
-   Total data transferred

------------------------------------------------------------------------

# Future Metrics

Possible future metrics include:

    deckhand_rsync_errors_total

These may be added as Deckhand evolves.
