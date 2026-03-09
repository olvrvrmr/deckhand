package docker

import (
	"strconv"
	"strings"
)

const (
	LabelEnable   = "deckhand.enable"
	LabelStop     = "deckhand.stop"
	LabelPath     = "deckhand.path"    // single host path to sync
	LabelExclude  = "deckhand.exclude" // comma-separated rsync exclude patterns
	LabelPreExec  = "deckhand.pre-exec"
	LabelPriority = "deckhand.priority" // stop order: lower = stopped first
)

type ContainerMeta struct {
	ID       string
	Name     string
	Stop     bool
	Path     string
	Excludes []string
	PreExec  string
	Priority int
}

func parseLabels(id, name string, labels map[string]string) *ContainerMeta {
	if labels[LabelEnable] != "true" {
		return nil
	}
	m := &ContainerMeta{
		ID:      id,
		Name:    name,
		Stop:    labels[LabelStop] == "true",
		Path:    labels[LabelPath],
		PreExec: labels[LabelPreExec],
	}
	if e := labels[LabelExclude]; e != "" {
		for _, pattern := range strings.Split(e, ",") {
			if pattern = strings.TrimSpace(pattern); pattern != "" {
				m.Excludes = append(m.Excludes, pattern)
			}
		}
	}
	if prio := labels[LabelPriority]; prio != "" {
		if i, err := strconv.Atoi(prio); err == nil {
			m.Priority = i
		}
	}
	return m
}
