package report

import (
	"time"
)

// ReportLog defines the Report log entity
type ReportLog struct {
	SessionID       string    `json:"sessionId"`
	Platform        string    `json:"platform"`
	ConnectUsername string    `json:"connectUsername"`
	EventType       string    `json:"eventType"`
	EventStartedAt  time.Time `json:"eventStartedAt"`
	EventEndedAt    time.Time `json:"eventEndedAt"`
	Duration        int       `json:"duration"` // in seconds
}
