// Copyright (c) 2026 Rasel Hossen
// See LICENSE for usage terms.

package workers

import (
	"context"
	"fmt"

	"github.com/rasel9t6/the-go-engineer/11-flagship/01-opslane/internal/events"
)

type Notification struct {
	TenantID int64
	EventID  string
	Channel  string
	Subject  string
	Body     string
}

type NotificationSink interface {
	Send(ctx context.Context, notification Notification) error
}

type NotificationWorker struct {
	Sink NotificationSink
}

func (w NotificationWorker) Handle(ctx context.Context, event events.Event) error {
	if event.Type != events.TypeNotificationRequested {
		return nil
	}
	if w.Sink == nil {
		return fmt.Errorf("notification sink is not configured")
	}
	if event.TenantID <= 0 {
		return events.ErrInvalidEvent
	}

	notification := Notification{
		TenantID: event.TenantID,
		EventID:  event.ID,
		Channel:  event.Metadata["channel"],
		Subject:  event.Metadata["subject"],
		Body:     event.Metadata["body"],
	}
	if notification.Channel == "" || notification.Subject == "" {
		return events.ErrInvalidEvent
	}

	if err := w.Sink.Send(ctx, notification); err != nil {
		return fmt.Errorf("send notification: %w", err)
	}

	return nil
}
