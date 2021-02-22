package config

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

var amqpConn *amqp.Connection

func (config AppConfig) SetupAMQPConnection() {
	// Open up messaging queue connection.
	conn, err := amqp.Dial(viper.GetString("MQ_DSN"))
	if err != nil {
		panic(fmt.Sprintf(
			"failed to connect to message broker, cause: %s",
			err.Error(),
		))
	}

	setAMQPConnection(conn)
}

func setAMQPConnection(currentAMQPConnection *amqp.Connection) {
	amqpConn = currentAMQPConnection
}

func (config AppConfig) GetAMQPConnection() *amqp.Connection {
	return amqpConn
}