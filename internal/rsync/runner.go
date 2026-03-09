package rsync

import (
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"strings"
)

type Runner struct {
	SSHKeyPath string
	ExtraArgs  []string
	DryRun     bool
}

func (r *Runner) Sync(src, dst string, excludes []string) error {
	args := []string{
		"-avz",
		"--delete",
		"--mkpath",
		"-e", fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no", r.SSHKeyPath),
	}
	for _, pattern := range excludes {
		args = append(args, "--exclude="+pattern)
	}
	if r.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, r.ExtraArgs...)

	if !strings.HasSuffix(src, "/") {
		src += "/"
	}
	args = append(args, src, dst)

	slog.Info("rsync", "src", src, "dst", dst, "dry_run", r.DryRun)
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rsync %s → %s: %w", src, dst, err)
	}
	return nil
}
