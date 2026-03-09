package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// BackupsTotal counts backup attempts per container and status (success/failure)
	BackupsTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "deckhand_backups_total",
		Help: "Total number of backup attempts per container.",
	}, []string{"container", "status"})

	// BackupFailuresTotal counts failed backups per container
	BackupFailuresTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "deckhand_backup_failures_total",
		Help: "Total number of failed backups per container.",
	}, []string{"container"})

	// BackupDuration tracks backup duration per container as a histogram
	BackupDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "deckhand_backup_duration_seconds",
		Help:    "Duration of each backup execution per container.",
		Buckets: prometheus.DefBuckets,
	}, []string{"container"})

	// LastBackupTimestamp records the unix timestamp of the last successful backup per container
	LastBackupTimestamp = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "deckhand_last_backup_timestamp",
		Help: "Unix timestamp of the last successful backup per container.",
	}, []string{"container"})

	// BytesTransferredTotal counts total bytes transferred per container
	BytesTransferredTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "deckhand_bytes_transferred_total",
		Help: "Total bytes transferred by rsync per container.",
	}, []string{"container"})
)

func init() {
	prometheus.MustRegister(
		BackupsTotal,
		BackupFailuresTotal,
		BackupDuration,
		LastBackupTimestamp,
		BytesTransferredTotal,
	)
}

func Serve(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(addr, nil)
}
