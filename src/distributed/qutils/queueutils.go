// Package qutils is responsible for the utilities
// needed for communication with the Message Broker.
// Need amqp package to interact with the
// Message Broker.
package qutils

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

const SensorListQueue = "SensorList"

// GetChannel is responsible for instantiating the connection
// and channel to communicate with the Message Broker.
// Not adding the get queue here in public function get channel
// as there can be multiple queues instantiated
func GetChannel(url string) (*amqp.Connection, *amqp.Channel) {
	//Instantiating connection
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to establish connection to Message Broker")

	//Instantiatinf channel on opened connection
	ch, err := conn.Channel()
	failOnError(err, "Failed to get channel for connection")

	return conn, ch
}

// GetQueue is responsible for declaring a new Queue
// using a Channel
// @param name | Routing key for the queue
// @param ch | Channel to declare the Queue
func GetQueue(name string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		name,  //name string,
		false, //durable bool,
		false, //autoDelete bool,
		false, //exclusive bool,
		false, //noWait bool,
		nil)   //args amqp.Table | Using the default exchange, no need of any other configuration information

	failOnError(err, "Failed to declare queue")

	return &q
}

//Responsible for logging the error and crashing the application
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
