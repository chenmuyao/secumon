package monitor

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func AccessLogExchangeSetup(conn *amqp.Connection, exchangeName string) error {
	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	return channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
}
