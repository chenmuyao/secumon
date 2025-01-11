package monitor

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/chenmuyao/secumon/internal/domain"
	amqp "github.com/rabbitmq/amqp091-go"
)

type LogMonitorConsumer interface {
	Consume(ctx context.Context, log domain.AccessLog) error
}

type RabbitMQLogMonitorConsumer struct {
	conn           *amqp.Connection
	ch             *amqp.Channel
	consumeTimeout time.Duration
}

// Consume implements LogMonitorConsumer.
func (r *RabbitMQLogMonitorConsumer) Consume(ctx context.Context, log domain.AccessLog) error {
	slog.Debug("Consuming log", slog.Any("log", log))
	return nil
}

func (r *RabbitMQLogMonitorConsumer) StartConsumer(
	exchangeName string,
	queueName string,
) error {
	var err error
	r.ch, err = r.conn.Channel()
	if err != nil {
		return err
	}

	deliveries, err := r.ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	go r.handle(deliveries)

	return nil
}

func (r *RabbitMQLogMonitorConsumer) handle(deliveries <-chan amqp.Delivery) {
	for d := range deliveries {
		var log domain.AccessLog
		err := json.Unmarshal(d.Body, &log)
		if err != nil {
			slog.Error("consumer unmarshal error", slog.Any("err", err))
			continue
		}
		ctx, cancel := context.WithTimeout(context.Background(), r.consumeTimeout)
		err = r.Consume(ctx, log)
		if err != nil {
			slog.Error("failed to consume", slog.Any("log", log), slog.Any("err", err))
			// just drop
			err = d.Nack(false, false)
			if err != nil {
				slog.Error("failed to nack", slog.Any("log", log), slog.Any("err", err))
			}
			cancel()
			continue
		}
		err = d.Ack(false)
		if err != nil {
			slog.Error("failed to ack", slog.Any("log", log), slog.Any("err", err))
		}
		cancel()
	}
}

func NewRabbitMQLogMonitorConsumer(conn *amqp.Connection) *RabbitMQLogMonitorConsumer {
	return &RabbitMQLogMonitorConsumer{
		conn:           conn,
		consumeTimeout: time.Second,
	}
}
