package events

import (
	"bytes"
	"fmt"
	"github.com/aasumitro/svc-telegram-notify/config"
	"github.com/aasumitro/svc-telegram-notify/src/helpers"
	"github.com/aasumitro/svc-telegram-notify/src/listeners"
	"github.com/spf13/viper"
	"log"
	"time"
)

// MessagingEvent represent the data-struct for configuration
type MessagingEvent struct {
	config       *config.AppConfig
}

func (event MessagingEvent) ListenToRabbitMQ() {
	// create messaging queue channel
	channel, err := event.config.GetAMQPConnection().Channel()
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to open concurrent server channel, cause: %s",
			err.Error(),
		))
	}

	// defer the close till after the main function has finished
	defer func() {
		err := channel.Close()
		if err != nil {
			panic(err)
		}
	}()

	// QueueDeclare declares a queue to hold messages and deliver to consumers.
	// Declaring creates a queue if it doesn't already exist, or ensures that an
	// existing queue matches the same parameters.
	queue, err := channel.QueueDeclare(
		viper.GetString("MQ_QUEUE_NAME"),
		viper.GetBool("MQ_QUEUE_DURABLE"),
		viper.GetBool("MQ_QUEUE_AUTO_DELETE"),
		viper.GetBool("MQ_QUEUE_EXCLUSIVE"),
		viper.GetBool("MQ_QUEUE_NO_WAIT"),nil)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to declare queue, cause: %s",
			err.Error(),
		))
	}

	// Qos controls how many messages or how many bytes the server will try to keep on
	// the network for consumers before receiving delivery Acknowledged. The intent of Qos is
	// to make sure the network buffers stay full between the server and client.
	err = channel.Qos(
		viper.GetInt("MQ_QOS_PREFETCH_COUNT"),
		viper.GetInt("MQ_QOS_PREFETCH_SIZE"),
		viper.GetBool("MQ_QOS_GLOBAL"))
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to set QoS, cause: %s",
			err.Error(),
		))
	}

	// Begin receiving on the returned chan Delivery.
	msg, err := channel.Consume(queue.Name, "",false,
		viper.GetBool("MQ_QUEUE_EXCLUSIVE"),
		false,
		viper.GetBool("MQ_QUEUE_NO_WAIT"),nil)
	if err != nil {
		panic(fmt.Sprintf(
			"Failed to register a consumer, cause: %s",
			err.Error(),
		))
	}

	// create an unbuffered channel for bool types.
	// Type is not important but we have to give one anyway.
	forever := make(chan bool)

	// fire up a goroutine that hooks onto message channel and reads
	// anything that pops into it. This essentially is a thread of
	// execution within the main thread. message is a channel constructed by
	// previous code.
	go func() {
		for delivery := range msg {
			// show log if new message is received
			fmt.Println(fmt.Sprintf("Received a message: %s", delivery.Body))
			fmt.Println("=====================================================")

			// make it happen
			event.handleEvent(delivery.Body)

			// -----------
			dotCount := bytes.Count(delivery.Body, []byte("."))
			t := time.Duration(dotCount)
			time.Sleep(t * time.Second)
			// -----------

			// Ack delegates an acknowledgement through the Acknowledged interface that the
			// client or server has finished work on a delivery.
			err := delivery.Ack(false)
			if err != nil {
				panic(fmt.Sprintf(
					"Failed to delegates an acknowledgement, cause: %s",
					err.Error(),
				))
			}
		}
	}()

	log.Printf(" [*] Start listening form messages broker (RabbitMQ)")
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	fmt.Println("=====================================================")

	// We need to block the main thread so that the above thread stays
	// on reading from msg channel. To do that just try to read in from
	// the forever channel. As long as no one writes to it we will wait here.
	// Since we are the only ones that know of it it is guaranteed that
	// nothing gets written in it. We could also do a busy wait here but
	// that would waste CPU cycles for no good reason.
	<-forever
}

func (event MessagingEvent) handleEvent(deliveryMessage []byte) {
	msg, err := helpers.Deserialize(deliveryMessage)
	if err != nil {
		fmt.Println(fmt.Sprintf(
			"Failed to deserialize delivery message, cause: %s",
			err.Error(),
		))
		return
	}

	listeners.InitMessagingListener(
		fmt.Sprint(msg["chat_id"]),
		fmt.Sprint(msg["message"]),
		event.config,
	).SendNotify()
}

// InitMessagingEvent initialize the app configuration
func InitMessagingEvent(
	applicationConfig	*config.AppConfig,
) *MessagingEvent {
	return &MessagingEvent{
		config: applicationConfig,
	}
}