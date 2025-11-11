package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kuahbanyak/go-crud/internal/adapters/external/email"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/logger"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/messaging/events"
	"github.com/kuahbanyak/go-crud/internal/infrastructure/messaging/rabbitmq"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		logger.Info("No .env file found, using system environment variables")
	}

	// RabbitMQ configuration
	rabbitConfig := rabbitmq.Config{
		Host:     getEnv("RABBITMQ_HOST", "localhost"),
		Port:     getEnv("RABBITMQ_PORT", "5672"),
		User:     getEnv("RABBITMQ_USER", "admin"),
		Password: getEnv("RABBITMQ_PASS", "password"),
		Vhost:    getEnv("RABBITMQ_VHOST", "/"),
	}

	// Connect to RabbitMQ
	conn, err := rabbitmq.NewConnection(rabbitConfig)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	defer conn.Close()

	// Setup infrastructure (exchanges, queues)
	if err := conn.SetupInfrastructure(); err != nil {
		log.Fatal("Failed to setup RabbitMQ infrastructure:", err)
	}

	// Initialize email service
	emailService := email.NewSMTPClient(
		getEnv("SMTP_HOST", "smtp.gmail.com"),
		getEnv("SMTP_PORT", "587"),
		getEnv("SMTP_USER", ""),
		getEnv("SMTP_PASS", ""),
	)

	// Start consuming messages
	logger.Info("Notification worker started. Waiting for messages...")

	// Email notifications consumer
	go consumeEmailNotifications(conn, emailService)

	// SMS notifications consumer (placeholder)
	go consumeSMSNotifications(conn)

	// Push notifications consumer (placeholder)
	go consumePushNotifications(conn)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down notification worker...")
}

func consumeEmailNotifications(conn *rabbitmq.Connection, emailService *email.SMTPClient) {
	msgs, err := conn.Consume("notifications.email", "email-worker")
	if err != nil {
		log.Fatal("Failed to start consuming email notifications:", err)
	}

	for msg := range msgs {
		if err := handleEmailNotification(msg, emailService); err != nil {
			logger.Error("Failed to handle email notification", err)

			// Retry logic
			if msg.Headers == nil {
				msg.Headers = make(amqp.Table)
			}

			retryCount, _ := msg.Headers["x-retry-count"].(int32)
			if retryCount < 3 {
				msg.Headers["x-retry-count"] = retryCount + 1
				msg.Nack(false, true) // Requeue
			} else {
				msg.Nack(false, false) // Send to DLQ
			}
		} else {
			msg.Ack(false)
		}
	}
}

func handleEmailNotification(msg amqp.Delivery, emailService *email.SMTPClient) error {
	var event events.EmailNotificationEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	logger.Info("Processing email notification", "to:", event.To)

	body := generateEmailBody(&event)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	_ = ctx // Use context if needed

	if err := emailService.SendEmail(event.To, event.Subject, body); err != nil {
		return err
	}

	logger.Info("Email sent successfully to", event.To)
	return nil
}

func generateEmailBody(event *events.EmailNotificationEvent) string {
	if event.Body != "" {
		return event.Body
	}

	switch event.Template {
	case "queue_assigned":
		return formatQueueAssignedEmail(event.TemplateData)
	case "service_started":
		return formatServiceStartedEmail(event.TemplateData)
	case "issue_discovered":
		return formatIssueDiscoveredEmail(event.TemplateData)
	case "approval_needed":
		return formatApprovalNeededEmail(event.TemplateData)
	case "service_completed":
		return formatServiceCompletedEmail(event.TemplateData)
	default:
		return "No template specified"
	}
}

func formatQueueAssignedEmail(data map[string]interface{}) string {
	return fmt.Sprintf("Dear %v,\n\nYour service booking confirmed!\nQueue #%v\nDate: %v\nType: %v\nVehicle: %v %v (%v)\n\nBest regards",
		data["customer_name"], data["queue_number"], data["service_date"], data["service_type"],
		data["vehicle_brand"], data["vehicle_model"], data["license_plate"])
}

func formatServiceStartedEmail(data map[string]interface{}) string {
	return fmt.Sprintf("Dear %v,\n\nYour service has started!\nQueue #%v\nType: %v\nVehicle: %v\n\nBest regards",
		data["customer_name"], data["queue_number"], data["service_type"], data["vehicle"])
}

func formatIssueDiscoveredEmail(data map[string]interface{}) string {
	return fmt.Sprintf("Dear %v,\n\nIssue discovered by %v:\n%v - %v\n%v\nPriority: %v\nCost: $%.2f\n\nBest regards",
		data["customer_name"], data["mechanic_name"], data["category"], data["item_name"],
		data["description"], data["priority"], data["estimated_cost"])
}

func formatApprovalNeededEmail(data map[string]interface{}) string {
	return fmt.Sprintf("Dear %v,\n\n%v items need approval\nTotal: $%.2f\nApprove: %v\n\nBest regards",
		data["customer_name"], data["item_count"], data["total_cost"], data["approval_url"])
}

func formatServiceCompletedEmail(data map[string]interface{}) string {
	return fmt.Sprintf("Dear %v,\n\nService complete!\nQueue #%v\nType: %v\nCompleted: %v\nCost: $%.2f\n\nThank you!",
		data["customer_name"], data["queue_number"], data["service_type"],
		data["completed_at"], data["total_cost"])
}

func consumeSMSNotifications(conn *rabbitmq.Connection) {
	msgs, err := conn.Consume("notifications.sms", "sms-worker")
	if err != nil {
		log.Fatal("Failed to start consuming SMS:", err)
	}

	for msg := range msgs {
		var event events.SMSNotificationEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			logger.Error("Failed to unmarshal SMS", err)
			msg.Nack(false, false)
			continue
		}

		logger.Info("SMS notification", "to:", event.To)
		msg.Ack(false)
	}
}

func consumePushNotifications(conn *rabbitmq.Connection) {
	msgs, err := conn.Consume("notifications.push", "push-worker")
	if err != nil {
		log.Fatal("Failed to start consuming push:", err)
	}

	for msg := range msgs {
		var event events.PushNotificationEvent
		if err := json.Unmarshal(msg.Body, &event); err != nil {
			logger.Error("Failed to unmarshal push", err)
			msg.Nack(false, false)
			continue
		}

		logger.Info("Push notification", "user:", event.UserID)
		msg.Ack(false)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
