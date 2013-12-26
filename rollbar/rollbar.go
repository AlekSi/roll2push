package rollbar

import (
	"encoding/json"
)

type Event struct {
	EventName string `json:"event_name"`
	Data      map[string]json.RawMessage
}

// event_name "*_item"
type ItemEvent struct {
	Environment    string
	Title          string
	LastOccurrence LastOccurrence `json:"last_occurrence"`
}

type LastOccurrence struct {
	Level     string
	Timestamp int64
}
