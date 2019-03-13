package coordinator

import (
	"bytes"
	"distributed/dto"
	"distributed/qutils"
	"encoding/gob"
	"fmt"

	"github.com/streadway/amqp"
)

// url to rabbit's end point so that
// connections can be established to it.
const url = "amqp://guest:guest@localhost:5672"

// QueueListener contains the logic that discovers the data queues,
// recevices the messages and eventually transales them to events
// in an EventAggregator
type QueueListener struct {
	conn *amqp.Connection // for getting messages
	ch   *amqp.Channel    // for getting messages

	// to prevent registering a seneor twice, to close of listener if associated sensor goes offline.
	// map that points to the Delivery objects
	sources map[string]<-chan amqp.Delivery //registry of all the sources that the coordinator is listening on
}

// NewQueueListener is a constructor function
// that ensures that the QueueListener is
// properly initialized
func NewQueueListener() *QueueListener {
	//instantiating new object of queuelistener
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery),
	}

	//populating the Connection and Channel fields
	ql.conn, ql.ch = qutils.GetChannel(url)

	return &ql
}

// ListenForNewSource is responsible for
// letting the QueueListener discover new sensors
func (ql *QueueListener) ListenForNewSource() {
	q := qutils.GetQueue("", ql.ch) // blank name for queue gives a random (unique) name to the queue (no conflicts when multiple coordinators are running)

	// By default the queue generated is bound to the default exchange.
	// the sensors publish to a fanout exchange, therefore, q needs to
	// rebind to that one.
	ql.ch.QueueBind(
		q.Name,       // name string,
		"",           // key string,
		"amq.fanout", // exchange string,
		false,        // noWait bool,
		nil)          // args amqp.Table

	// Receiver for consuming the messages
	msgs, _ := ql.ch.Consume(
		q.Name, //name of the queue bound to the fanout exchange
		"",
		true,
		false,
		false,
		false,
		nil)

	// Channel in place, waiting for the messages on the msgs channel
	for msg := range msgs {
		// new message mesans a new sensor is online
		// and is ready to send the readings.
		// usind consume method to get access to that queue.
		sourceChan, _ := ql.ch.Consume(
			string(msg.Body), //name is in the msg body for the data queue.
			"",
			true,
			false,
			false,
			false,
			nil)
		// Sending data in a default exchange, this is a direct exchange
		// and will only deliver a message to a single receiver. That means
		// when multiple coordinatos are registered, they share access to the queues
		// when this happens, rabbitmq will take turns delivering to each registers
		// receiver in turn. This lets us scale the coordinators as the system grows
		// without affecting the rest of the system.

		//checking if new message has already been registered
		if ql.sources[string(msg.Body)] == nil {
			ql.sources[string(msg.Body)] = sourceChan

			go ql.AddListener(sourceChan)
		}
	}
}

// AddListener is responsible
func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) {
	// waiting for messages from the channel
	for msg := range msgs {
		// publish events for the downstream consumers
		// convert binary data to workable data
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.SensorMessage)
		d.Decode(sd)

		fmt.Printf("Received message: %v\n", sd)
	}
}
