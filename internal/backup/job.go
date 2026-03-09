package backup

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/obornet/deckhand/internal/config"
	"github.com/obornet/deckhand/internal/docker"
	"github.com/obornet/deckhand/internal/notify"
	"github.com/obornet/deckhand/internal/rsync"
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

	if err := j.run(ctx); err != nil {
		slog.Error("backup failed", "error", err, "duration", time.Since(start))
		notify.Send(j.cfg.NotifyURL, false, err.Error())
		return
	}

	slog.Info("backup completed", "duration", time.Since(start))
	notify.Send(j.cfg.NotifyURL, true, "")
}

func (j *Job) run(ctx context.Context) error {
	containers, err := j.docker.GetBackupContainers(ctx)
	if err != nil {
		return fmt.Errorf("discover containers: %w", err)
	}

	if len(containers) == 0 {
		slog.Info("no containers with deckhand.enable=true found")
		return nil
	}

	// stop containers that need it (priority order)
	var stopped []string
	for _, c := range containers {
		if c.Stop {
			slog.Info("stopping container", "name", c.Name)
			if err := j.docker.Stop(ctx, c.ID); err != nil {
				return fmt.Errorf("stop %s: %w", c.Name, err)
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

	// run pre-exec hooks on running containers (stop=false only, or after stop would be pointless)
	for _, c := range containers {
		if c.PreExec != "" && !c.Stop {
			slog.Info("running pre-exec", "container", c.Name, "cmd", c.PreExec)
			if err := j.docker.Exec(ctx, c.ID, c.PreExec); err != nil {
				return fmt.Errorf("pre-exec %s: %w", c.Name, err)
			}
		}
	}

	// sync paths for each container
	for _, c := range containers {
		for _, path := range c.Paths {
			dst := fmt.Sprintf("%s/%s%s", j.cfg.Destination, c.Name, path)
			if err := j.rsync.Sync(path, dst); err != nil {
				return err
			}
		}
	}

	return nil
}
