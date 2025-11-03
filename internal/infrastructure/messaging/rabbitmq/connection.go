package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Vhost    string
}

type Connection struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  Config
}

func NewConnection(config Config) (*Connection, error) {
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s%s",
		config.User, config.Password, config.Host, config.Port, config.Vhost)

	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = channel.Qos(10, 0, false)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	log.Println("Successfully connected to RabbitMQ")

	return &Connection{conn: conn, channel: channel, config: config}, nil
}

func (c *Connection) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Connection) Channel() *amqp.Channel {
	return c.channel
}

func (c *Connection) DeclareExchange(name, kind string, durable bool) error {
	return c.channel.ExchangeDeclare(name, kind, durable, false, false, false, nil)
}

func (c *Connection) DeclareQueue(name string, durable bool, dlx string) (amqp.Queue, error) {
	args := amqp.Table{}
	if dlx != "" {
		args["x-dead-letter-exchange"] = dlx
	}
	return c.channel.QueueDeclare(name, durable, false, false, false, args)
}

func (c *Connection) BindQueue(queueName, exchangeName, routingKey string) error {
	return c.channel.QueueBind(queueName, routingKey, exchangeName, false, nil)
}

func (c *Connection) Publish(ctx context.Context, exchange, routingKey string, body []byte) error {
	return c.channel.PublishWithContext(ctx, exchange, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Timestamp:    time.Now(),
		Body:         body,
	})
}

func (c *Connection) PublishWithRetry(ctx context.Context, exchange, routingKey string, body []byte, maxRetries int) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		err = c.Publish(ctx, exchange, routingKey, body)
		if err == nil {
			return nil
		}
		time.Sleep(time.Second * time.Duration(i+1))
	}
	return fmt.Errorf("failed to publish after %d attempts: %w", maxRetries, err)
}

func (c *Connection) Consume(queueName, consumerTag string) (<-chan amqp.Delivery, error) {
	return c.channel.Consume(queueName, consumerTag, false, false, false, false, nil)
}

func (c *Connection) SetupInfrastructure() error {
	if err := c.DeclareExchange("car-maintenance", "topic", true); err != nil {
		return err
	}
	if err := c.DeclareExchange("car-maintenance.dlx", "topic", true); err != nil {
		return err
	}

	queues := []struct {
		name       string
		routingKey string
	}{
		{"notifications.email", "notification.email"},
		{"notifications.sms", "notification.sms"},
		{"notifications.push", "notification.push"},
		{"events.queue-status", "event.queue.*"},
		{"events.service-status", "event.service.*"},
		{"events.approval", "event.approval.*"},
		{"payments.process", "payment.process"},
		{"audit.log", "audit.*"},
	}

	for _, q := range queues {
		queue, err := c.DeclareQueue(q.name, true, "car-maintenance.dlx")
		if err != nil {
			return err
		}
		if err := c.BindQueue(queue.Name, "car-maintenance", q.routingKey); err != nil {
			return err
		}
	}

	dlq, err := c.DeclareQueue("dead-letter-queue", true, "")
	if err != nil {
		return err
	}
	if err := c.BindQueue(dlq.Name, "car-maintenance.dlx", "#"); err != nil {
		return err
	}

	log.Println("RabbitMQ infrastructure setup completed")
	return nil
}
