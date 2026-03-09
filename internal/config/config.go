package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	CronSchedule   string
	Destination    string
	SSHKeyPath     string
	ExtraRsyncArgs []string
	NotifyURL      string
	DryRun         bool
	RunOnce        bool
}

func Load() *Config {
	c := &Config{
		CronSchedule: getEnv("BACKUP_CRON", "0 2 * * *"),
		Destination:  getEnv("BACKUP_DESTINATION", ""),
		SSHKeyPath:   getEnv("BACKUP_SSH_KEY", "/keys/id_rsa"),
		NotifyURL:    getEnv("BACKUP_NOTIFY_URL", ""),
		DryRun:       getEnvBool("BACKUP_DRY_RUN", false),
		RunOnce:      getEnvBool("BACKUP_RUN_ONCE", false),
	}
	if extra := getEnv("BACKUP_RSYNC_ARGS", ""); extra != "" {
		c.ExtraRsyncArgs = strings.Fields(extra)
	}
	return c
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
