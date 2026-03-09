package main

import (
	"log/slog"
	"os"

	"github.com/olvrvrmr/deckhand/internal/backup"
	"github.com/olvrvrmr/deckhand/internal/config"
	"github.com/olvrvrmr/deckhand/internal/docker"
	"github.com/olvrvrmr/deckhand/internal/metrics"
	"github.com/robfig/cron/v3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	cfg := config.Load()
	if cfg.Destination == "" {
		slog.Error("BACKUP_DESTINATION is required")
		os.Exit(1)
	}

	dockerClient, err := docker.New()
	if err != nil {
		slog.Error("failed to connect to Docker", "error", err)
		os.Exit(1)
	}

	job := backup.New(cfg, dockerClient)

	if !cfg.RunOnce {
		metrics.Serve(cfg.MetricsAddr)
		slog.Info("metrics server started", "addr", cfg.MetricsAddr)
	}

	if cfg.RunOnce {
		slog.Info("running once")
		job.Run()
		return
	}

	c := cron.New()
	if _, err := c.AddFunc(cfg.CronSchedule, job.Run); err != nil {
		slog.Error("invalid cron schedule", "schedule", cfg.CronSchedule, "error", err)
		os.Exit(1)
	}
	c.Start()
	slog.Info("scheduler started", "cron", cfg.CronSchedule)
	select {}
}
