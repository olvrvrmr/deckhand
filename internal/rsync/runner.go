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

func (r *Runner) Sync(src, dst string) error {
	args := []string{
		"-avz",
		"--delete",
		"-e", fmt.Sprintf("ssh -i %s -o StrictHostKeyChecking=no", r.SSHKeyPath),
	}
	if r.DryRun {
		args = append(args, "--dry-run")
	}
	args = append(args, r.ExtraArgs...)

	// trailing slash on src = sync contents, not the directory itself
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
