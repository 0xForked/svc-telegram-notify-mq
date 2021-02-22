## Notify Service

### What's this?
[TelegramBot Webhook and Worker](https://core.telegram.org/bots) implementation send notification as chat with [golang](https://golang.org/) (Go Programming Language).

### How it's works?

This service subscribe and handle an incoming queue's event from [Message Broker](https://medium.com/@acep.abdurohman90/mengapa-menggunakan-message-broker-c17453cb225e)
([RabbitMQ](https://www.rabbitmq.com/)) via [amqp](https://www.amqp.org/) connection protocol and sent it to users via
telegeram bot using [gotgbot](https://github.com/PaulSonOfLars/gotgbot) library for golang.

<hr>
<p align="center">
The part of microservices architecture!
</p>