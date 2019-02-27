/**
 * Sensor application
 *
 * Stand alone application for mocking
 * a real time sensor.
 */
package main

import (
	"flag"
	"log"
	"math/rand"
	"strconv"
	"time"
)

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

/**
 * Periodically generate data points that
 * drift between the maximum and minimum
 * values that are provided to the sensor
 * with a bias towards putting the data points
 * near the average value.
 */
func main() {
	flag.Parse()

	//miliseconds per cycle
	//eg. 5 cycles / sec = 200 miliseconds / cycle
	dur, _ := time.ParseDuration(strconv.Itoa(1000/int(*freq)) + "ms")

	//channel that responds after the above created duration
	signal := time.Tick(dur)

	for range signal {
		//generating the values
		calcValue()
		log.Printf("Reading sent. Value: %v\n", value)
	}
}

/**
 * Responsible for calculating the data point.
 */
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
