package monitor

import (
	"net/http"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

var StatusList = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusUnauthorized,
	http.StatusInternalServerError,
}

func AccessLogMQSetup(conn *amqp.Connection, exchangeName string) (string, error) {
	channel, err := conn.Channel()
	if err != nil {
		return "", err
	}
	defer channel.Close()

	err = channel.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		return "", err
	}

	q, err := channel.QueueDeclare("", true, false, true, false, nil)

	for _, status := range StatusList {
		err = channel.QueueBind(q.Name, strconv.Itoa(status), exchangeName, false, nil)
		if err != nil {
			return "", err
		}
	}

	return q.Name, nil
}
