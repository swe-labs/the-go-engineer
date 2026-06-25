package main

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	ID        int             `json:"id"`
	Title     string          `json:"title"`
	Timestamp string          `json:"timestamp"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

func UnmarshalEvent(data []byte) (Event, error) {
	var e Event
	if err := json.Unmarshal(data, &e); err != nil {
		return Event{}, err
	}
	return e, nil
}

func main() {
	valid := []byte(`{
		"id": 1,
		"title": "deploy",
		"timestamp": "2025-06-01T10:00:00Z",
		"payload": {"env": "prod"}
	}`)
	e, err := UnmarshalEvent(valid)
	if err != nil {
		fmt.Println("Valid case error:", err)
	} else {
		fmt.Printf("Valid: ID=%d Title=%q TS=%q Payload=%s\n", e.ID, e.Title, e.Timestamp, string(e.Payload))
	}

	wrongType := []byte(`{"id": "not-a-number", "title": "x", "timestamp": "2025-01-01T00:00:00Z"}`)
	_, err = UnmarshalEvent(wrongType)
	fmt.Println("Wrong type error:", err)

	noTimestamp := []byte(`{"id": 2, "title": "test"}`)
	e2, err := UnmarshalEvent(noTimestamp)
	if err != nil {
		fmt.Println("No timestamp error:", err)
	} else {
		fmt.Printf("No timestamp: ID=%d Title=%q TS=%q\n", e2.ID, e2.Title, e2.Timestamp)
	}
}
