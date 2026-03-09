package rsync

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

var sentBytesRe = regexp.MustCompile(`sent ([\d,]+) bytes`)

type Result struct {
	BytesTransferred int64
}

type Runner struct {
	SSHKeyPath string
	ExtraArgs  []string
	DryRun     bool
}

func (r *Runner) Sync(src, dst string, excludes []string) (Result, error) {
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

	var buf bytes.Buffer
	cmd := exec.Command("rsync", args...)
	cmd.Stdout = io.MultiWriter(os.Stdout, &buf)
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return Result{}, fmt.Errorf("rsync %s → %s: %w", src, dst, err)
	}

	return Result{BytesTransferred: parseSentBytes(buf.String())}, nil
}

func parseSentBytes(output string) int64 {
	m := sentBytesRe.FindStringSubmatch(output)
	if len(m) < 2 {
		return 0
	}
	cleaned := strings.ReplaceAll(m[1], ",", "")
	n, err := strconv.ParseInt(cleaned, 10, 64)
	if err != nil {
		return 0
	}
	return n
}
