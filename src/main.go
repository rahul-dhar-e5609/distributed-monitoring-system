package main

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	go server()

	var a string
	fmt.Scanln(&a)
}

/**
 * Server for checking Rabbit MQ
 */
func server() {

	//Fetching the queue
	conn, ch, q := getQueue()

	//closig the connection and channel on defer
	defer conn.Close()
	defer ch.Close()

	//amqp library's Publishing struct acts as message
	//which is published via channel
	msg := amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte("Hello RabbitMQ"),
	}
	for {
		//Blank string for exchange means default exchange
		ch.Publish("", q.Name, false, false, msg)
	}
}

/**
 * Helper function responsible for returning a message queue.
 *
 * Using Default Exchange for Rabbit MQ. Therefore the name of
 * the queues happen to be the routing keys and thus it looks
 * as if there is no Exchange here.
 *
 * @return *amqp.Connection | Actual connection between the application and Rabbit MQ
 * @return *amqp.Channel | Communication between the application and Rabbit MQ
 * @return *amqp.Queue | Message queue that can be accessed using the channel
 */
func getQueue() (*amqp.Connection, *amqp.Channel, *amqp.Queue) {

	//creating connection
	conn, err := amqp.Dial("amqp://guest@rahul.dhar:5672")
	failOnError(err, "Failed to connect to RabbitMQ")

	//channel to communicate on the network
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	//Message Queue
	q, err := ch.QueueDeclare("hello",
		false, //durable bool,
		false, //autoDelete bool,
		false, //exclusivebool,
		false, //noWait bool,
		nil)   //args amqp.Table) // used for declaring the headers.
	failOnError(err, "Failed to declare a queue")

	//returning the connection, channel and queue
	return conn, ch, &q
}

/**
 * Helper function to fail on error and crash the application.amqp
 *
 * @param err | error object
 * @param msg | Failure Message
 */
func failOnError(err error, msg string) {
	if err != nil {
		//logging error
		log.Fatalf("%s: %s", msg, err)

		//crashing the application.
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
