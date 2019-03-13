// Sensor application
// Stand alone application for mocking
// a real time sensor.
package main

import (
	"bytes"
	"distributed/dto"
	"distributed/qutils"
	"encoding/gob"
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

//Broker's input listener
var url = "amqp://guest:guest@rahul.dhar:5672"

//flags to read the parameters in to config the sensors from the command line.

//unique name of each sensor
//will be used for the routing key for the messages
//therefore each value should be unique, otherwise
//sensor data will get mixed up
var name = flag.String("name", "sensor", "name of the sensor")

//How many data points generated per second
var freq = flag.Uint("freq", 5, "update frequency in cycles/sec")

//Max - min define the absolute limits of the measurements range
var max = flag.Float64("max", 5., "maximum values for geneerated readings")
var min = flag.Float64("min", 1., "minimum values for geneerated readings")

//Maximum allowable change
var stepSize = flag.Float64("step", 0.1, "maximum allowable change per measurement")

//rand.New needs a source to create the random number
//rand.NewSource needs a seed
var r = rand.New(rand.NewSource(time.Now().UnixNano()))

//initial value of the sensor
var value = r.Float64()*(*max-*min) + *min

//data point changes with every time step by a maximum of stepsize
//nominal value to bias the step so that it trends towards it.
var nom = (*max-*min)/2 + *min

// Periodically generate data points that
// drift between the maximum and minimum
// values that are provided to the sensor
// with a bias towards putting the data points
// near the average value.
func main() {
	flag.Parse()

	//Connection and channel for message broker
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	// Fetching the queue using the channel and routing key
	// The queue needs to be declared here even though we are not gonna
	// write to the queue directly. Rabbit MQ receives messages and
	// exchanges and uses their configuration and information within
	// the message to determine which queues to deliver to, so technically
	// only the routing key needs to be given when publishing.
	// However that isnt enough to ensure that a queue with that
	// routing key actually exists. By declaring the queue here, we can ensure that,
	dataQueue := qutils.GetQueue(*name, ch)

	// The queue on which the sensor module sends the name of the queue
	// whenever a new queue comes online, so that the coordinators able to
	// understand that they need to take messages from this new queue
	// @deprecated
	// was created to ensure that messages are received
	// this will now be the responsibility of the consumers
	// since they will now have to create thier own queue to listen on this exchange
	// sensorQueue := qutils.GetQueue(qutils.SensorListQueue, ch)

	// publishing presence of a new queue
	publishQueueName(ch)

	// creating a new queue
	discoveryQueue := qutils.GetQueue("", ch)

	// binding the above created queue too the sensor discovery package
	// so that it knows when the coordinators make a discovery request
	ch.QueueBind(
		discoveryQueue.Name,            //name
		"",                             // key
		qutils.SensorDiscoveryExchange, //exchange
		false,
		nil)

	//listen for discovery request
	go listenForDiscoveryRequests(discoveryQueue.Name, ch)

	//miliseconds per cycle
	//eg. 5 cycles / sec = 200 miliseconds / cycle
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	//channel that responds after the above created duration
	signal := time.Tick(dur)

	//buffer for encoding the message for transmission
	buf := new(bytes.Buffer)

	//to encode the message directly to the buffer
	enc := gob.NewEncoder(buf)

	for range signal {
		//generating the values
		calcValue()
		reading := dto.SensorMessage{
			Name:      *name,
			Value:     value,
			Timestamp: time.Now(),
		}
		//reseting buffer for new signal message
		//any previous data gets removed and buffers
		//pointer is reset to its initial position
		buf.Reset()

		// needs to be recreated every tiem it is used
		enc = gob.NewEncoder(buf)

		//encoding the reading
		enc.Encode(reading)

		msg := amqp.Publishing{
			Body: buf.Bytes(),
		}

		ch.Publish(
			"",             //exchange string, | Using default
			dataQueue.Name, //key string, | Routing key for the queue
			false,          //mandatory bool,
			false,          //immediate bool,
			msg)            //msg amqp.Publishing

		log.Printf("Reading sent. Value: %v\n", value)
	}
}

// listenForDiscoveryRequests is responsible for
// listening and receiving the discovery requests
// from the coordinators
func listenForDiscoveryRequests(name string, ch *amqp.Channel) {
	// receiving discovery requests
	// from the coordinators.
	msgs, _ := ch.Consume(name, "", true, false, false, false, nil)

	for range msgs {
		// publishing name of queue on
		// fanout exchange.
		publishQueueName(ch)
	}
}

// This func is responsible for publishing
// the presece of data queue when the
// sensor starts up
func publishQueueName(ch *amqp.Channel) {
	// message to sensor queue should be the name of the newly generated queue
	msg := amqp.Publishing{
		Body: []byte(*name),
	}
	ch.Publish(
		"amq.fanout", // "", changing from default ro fanout exchange
		"",           // sensorQueue.Name, changed to empty string as wont be needed ny fanout exchange
		false,
		false,
		msg)
}

// Responsible for calculating the data point.
func calcValue() {
	//maximum value the data point can decrease or increase
	var maxStep, minStep float64

	if value < nom {
		maxStep = *stepSize
		minStep = -1 * *stepSize * (value - *min) / (nom - *min)

	} else {
		maxStep = *stepSize * (*max - value) / (*max - nom)
		minStep = -1 * *stepSize
	}

	//random valur in the range and added it to the current value
	value += r.Float64()*(maxStep-minStep) + minStep

}
