package docker

import (
	"strconv"
	"strings"
)

const (
	LabelEnable   = "deckhand.enable"
	LabelStop     = "deckhand.stop"
	LabelPaths    = "deckhand.paths"    // comma-separated host paths to sync
	LabelPreExec  = "deckhand.pre-exec" // command to run inside container before sync
	LabelPriority = "deckhand.priority" // stop order: lower = stopped first
)

type ContainerMeta struct {
	ID       string
	Name     string
	Stop     bool
	Paths    []string
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
		PreExec: labels[LabelPreExec],
	}
	if p := labels[LabelPaths]; p != "" {
		for _, path := range strings.Split(p, ",") {
			if path = strings.TrimSpace(path); path != "" {
				m.Paths = append(m.Paths, path)
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
