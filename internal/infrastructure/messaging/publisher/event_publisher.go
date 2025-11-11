package publisher

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/messaging/events"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/messaging/rabbitmq"
)

type EventPublisher struct {
	conn *rabbitmq.Connection
}

func NewEventPublisher(conn *rabbitmq.Connection) *EventPublisher {
	return &EventPublisher{conn: conn}
}

func (p *EventPublisher) PublishQueueNumberAssigned(ctx context.Context, event *events.QueueNumberAssignedEvent) error {
	event.BaseEvent = events.BaseEvent{
		ID: uuid.New().String(), Type: events.EventQueueNumberAssigned,
		Timestamp: time.Now(), Source: "api",
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.PublishWithRetry(ctx, "car-maintenance", "event.queue.assigned", body, 3); err != nil {
		logger.Error("Failed to publish queue assigned event", err)
		return err
	}

	emailEvent := &events.EmailNotificationEvent{
		BaseEvent: events.BaseEvent{
			ID: uuid.New().String(), Type: events.EventNotificationEmail,
			Timestamp: time.Now(), Source: "api",
		},
		To:       event.CustomerEmail,
		Subject:  fmt.Sprintf("Booking Confirmed - Queue #%d", event.QueueNumber),
		Template: "queue_assigned",
		TemplateData: map[string]interface{}{
			"customer_name": event.CustomerName, "queue_number": event.QueueNumber,
			"service_date": event.ServiceDate.Format("Monday, January 2, 2006"),
			"service_type": event.ServiceType, "vehicle_brand": event.VehicleBrand,
			"vehicle_model": event.VehicleModel, "license_plate": event.LicensePlate,
		},
		Priority: "high",
	}

	return p.PublishEmailNotification(ctx, emailEvent)
}

func (p *EventPublisher) PublishServiceStarted(ctx context.Context, event *events.ServiceStartedEvent) error {
	event.BaseEvent = events.BaseEvent{
		ID: uuid.New().String(), Type: events.EventServiceStarted,
		Timestamp: time.Now(), Source: "api",
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.PublishWithRetry(ctx, "car-maintenance", "event.service.started", body, 3); err != nil {
		logger.Error("Failed to publish service started event", err)
		return err
	}

	emailEvent := &events.EmailNotificationEvent{
		BaseEvent: events.BaseEvent{
			ID: uuid.New().String(), Type: events.EventNotificationEmail,
			Timestamp: time.Now(), Source: "api",
		},
		To:       event.CustomerEmail,
		Subject:  "Your Service Has Started",
		Template: "service_started",
		TemplateData: map[string]interface{}{
			"customer_name": event.CustomerName, "queue_number": event.QueueNumber,
			"service_type": event.ServiceType,
			"vehicle":      fmt.Sprintf("%s %s", event.VehicleBrand, event.VehicleModel),
		},
		Priority: "high",
	}

	return p.PublishEmailNotification(ctx, emailEvent)
}

func (p *EventPublisher) PublishIssueDiscovered(ctx context.Context, event *events.IssueDiscoveredEvent) error {
	event.BaseEvent = events.BaseEvent{
		ID: uuid.New().String(), Type: events.EventIssueDiscovered,
		Timestamp: time.Now(), Source: "api",
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.PublishWithRetry(ctx, "car-maintenance", "event.approval.issue_discovered", body, 3); err != nil {
		logger.Error("Failed to publish issue discovered event", err)
		return err
	}

	if event.RequiresApproval {
		emailEvent := &events.EmailNotificationEvent{
			BaseEvent: events.BaseEvent{
				ID: uuid.New().String(), Type: events.EventNotificationEmail,
				Timestamp: time.Now(), Source: "api",
			},
			To:       event.CustomerEmail,
			Subject:  "⚠️ Additional Issue Found - Approval Needed",
			Template: "issue_discovered",
			TemplateData: map[string]interface{}{
				"customer_name": event.CustomerName, "mechanic_name": event.MechanicName,
				"category": event.Category, "item_name": event.ItemName,
				"description": event.Description, "priority": event.Priority,
				"estimated_cost": event.EstimatedCost, "image_url": event.ImageURL,
			},
			Priority: "high",
		}
		return p.PublishEmailNotification(ctx, emailEvent)
	}

	return nil
}

func (p *EventPublisher) PublishServiceCompleted(ctx context.Context, event *events.ServiceCompletedEvent) error {
	event.BaseEvent = events.BaseEvent{
		ID: uuid.New().String(), Type: events.EventServiceCompleted,
		Timestamp: time.Now(), Source: "api",
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	if err := p.conn.PublishWithRetry(ctx, "car-maintenance", "event.service.completed", body, 3); err != nil {
		logger.Error("Failed to publish service completed event", err)
		return err
	}

	emailEvent := &events.EmailNotificationEvent{
		BaseEvent: events.BaseEvent{
			ID: uuid.New().String(), Type: events.EventNotificationEmail,
			Timestamp: time.Now(), Source: "api",
		},
		To:       event.CustomerEmail,
		Subject:  "Service Completed - Thank You!",
		Template: "service_completed",
		TemplateData: map[string]interface{}{
			"customer_name": event.CustomerName, "queue_number": event.QueueNumber,
			"service_type": event.ServiceType, "completed_at": event.CompletedAt.Format("3:04 PM"),
			"total_cost": event.TotalCost, "items_completed": event.ItemsCompleted,
		},
		Priority: "normal",
	}

	return p.PublishEmailNotification(ctx, emailEvent)
}

func (p *EventPublisher) PublishEmailNotification(ctx context.Context, event *events.EmailNotificationEvent) error {
	if event.ID == "" {
		event.BaseEvent = events.BaseEvent{
			ID: uuid.New().String(), Type: events.EventNotificationEmail,
			Timestamp: time.Now(), Source: "api",
		}
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal email event: %w", err)
	}

	return p.conn.PublishWithRetry(ctx, "car-maintenance", "notification.email", body, 3)
}

func (p *EventPublisher) PublishSMSNotification(ctx context.Context, event *events.SMSNotificationEvent) error {
	if event.ID == "" {
		event.BaseEvent = events.BaseEvent{
			ID: uuid.New().String(), Type: events.EventNotificationSMS,
			Timestamp: time.Now(), Source: "api",
		}
	}

	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal SMS event: %w", err)
	}

	return p.conn.PublishWithRetry(ctx, "car-maintenance", "notification.sms", body, 3)
}
