package backup

import (
	"context"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/olvrvrmr/deckhand/internal/config"
	"github.com/olvrvrmr/deckhand/internal/docker"
	"github.com/olvrvrmr/deckhand/internal/metrics"
	"github.com/olvrvrmr/deckhand/internal/notify"
	"github.com/olvrvrmr/deckhand/internal/rsync"
)

type Job struct {
	cfg    *config.Config
	docker *docker.Client
	rsync  *rsync.Runner
}

func New(cfg *config.Config, d *docker.Client) *Job {
	return &Job{
		cfg:    cfg,
		docker: d,
		rsync: &rsync.Runner{
			SSHKeyPath: cfg.SSHKeyPath,
			ExtraArgs:  cfg.ExtraRsyncArgs,
			DryRun:     cfg.DryRun,
		},
	}
}

func (j *Job) Run() {
	ctx := context.Background()
	start := time.Now()
	slog.Info("backup started")

	count, err := j.run(ctx)
	duration := time.Since(start)
	metrics.BackupDuration.Set(duration.Seconds())

	if err != nil {
		slog.Error("backup failed", "error", err, "duration", duration)
		metrics.BackupFailure.Inc()
		notify.Send(j.cfg.NotifyURL, false, err.Error())
		return
	}

	metrics.BackupSuccess.Inc()
	metrics.BackupLastSuccess.SetToCurrentTime()
	metrics.ContainersBackedUp.Set(float64(count))
	slog.Info("backup completed", "duration", duration, "containers", count)
	notify.Send(j.cfg.NotifyURL, true, "")
}

func (j *Job) run(ctx context.Context) (int, error) {
	containers, err := j.docker.GetBackupContainers(ctx)
	if err != nil {
		return 0, fmt.Errorf("discover containers: %w", err)
	}

	if len(containers) == 0 {
		slog.Info("no containers with deckhand.enable=true found")
		return 0, nil
	}

	// stop containers that need it (priority order)
	var stopped []string
	for _, c := range containers {
		if c.Stop {
			slog.Info("stopping container", "name", c.Name)
			if err := j.docker.Stop(ctx, c.ID); err != nil {
				return 0, fmt.Errorf("stop %s: %w", c.Name, err)
			}
			stopped = append(stopped, c.ID)
		}
	}

	// always restart stopped containers, even on error
	defer func() {
		for _, id := range stopped {
			slog.Info("starting container", "id", id[:12])
			if err := j.docker.Start(ctx, id); err != nil {
				slog.Error("failed to restart container", "id", id[:12], "error", err)
			}
		}
	}()

	// run pre-exec hooks on running containers
	for _, c := range containers {
		if c.PreExec != "" && !c.Stop {
			slog.Info("running pre-exec", "container", c.Name, "cmd", c.PreExec)
			if err := j.docker.Exec(ctx, c.ID, c.PreExec); err != nil {
				return 0, fmt.Errorf("pre-exec %s: %w", c.Name, err)
			}
		}
	}

	// sync each container's path
	synced := 0
	for _, c := range containers {
		if c.Path == "" {
			slog.Warn("no path defined, skipping sync", "container", c.Name)
			continue
		}
		dst := fmt.Sprintf("%s/%s", j.cfg.Destination, filepath.Base(c.Path))
		if err := j.rsync.Sync(c.Path, dst, c.Excludes); err != nil {
			return synced, err
		}
		synced++
	}

	return synced, nil
}
