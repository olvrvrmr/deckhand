package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	BackupSuccess = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "deckhand_backup_success_total",
		Help: "Total number of successful backup runs.",
	})

	BackupFailure = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "deckhand_backup_failure_total",
		Help: "Total number of failed backup runs.",
	})

	BackupDuration = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "deckhand_backup_duration_seconds",
		Help: "Duration of the last backup run in seconds.",
	})

	BackupLastSuccess = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "deckhand_backup_last_success_timestamp",
		Help: "Unix timestamp of the last successful backup run.",
	})

	ContainersBackedUp = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "deckhand_containers_backed_up",
		Help: "Number of containers synced in the last backup run.",
	})
)

func init() {
	prometheus.MustRegister(
		BackupSuccess,
		BackupFailure,
		BackupDuration,
		BackupLastSuccess,
		ContainersBackedUp,
	)
}

func Serve(addr string) {
	http.Handle("/metrics", promhttp.Handler())
	go http.ListenAndServe(addr, nil)
}
