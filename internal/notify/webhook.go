package notify

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

type payload struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Time    string `json:"time"`
}

func Send(url string, success bool, message string) {
	if url == "" {
		return
	}
	p := payload{
		Success: success,
		Message: message,
		Time:    time.Now().Format(time.RFC3339),
	}
	body, _ := json.Marshal(p)
	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Warn("notify failed", "error", err)
		return
	}
	defer resp.Body.Close()
}
