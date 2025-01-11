package monitor

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/chenmuyao/secumon/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LogMonitorPublisher interface {
	Publish(ctx context.Context, log domain.AccessLog) error
}

type RabbitMQLogMonitorPublisher struct {
	conn         *amqp.Connection
	exchangeName string
	immediate    bool
}

func NewRabbitMQLogMonitorPublisher(conn *amqp.Connection, topicName string) LogMonitorPublisher {
	return &RabbitMQLogMonitorPublisher{
		conn:         conn,
		exchangeName: topicName,
	}
}

// Publish implements LogMonitorPublisher.
// NOTE: This is the simplest version of Publisher. Several improvements could be done on demand:
// - Use a channel pool to no creating a channel at each publishment
// - Use a better json pkg for marshalling
func (r *RabbitMQLogMonitorPublisher) Publish(ctx context.Context, log domain.AccessLog) error {
	ch, err := r.conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	logBytes, err := json.Marshal(log)
	if err != nil {
		return err
	}

	return ch.PublishWithContext(
		ctx,
		r.exchangeName,
		strconv.Itoa(log.StatusCode),
		// The queue is predefined, so it must exist.
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logBytes,
		},
	)
}
