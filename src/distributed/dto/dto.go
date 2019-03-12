// Package dto This holds the definition
// for any types that we want to send through
// the message broker.
// Data Transfer Objects
package dto

import (
	"encoding/gob"
	"time"
)

//SensorMessage carries information to the coordinator
type SensorMessage struct {
	//name of the sensor, so that we
	//know where the data is coming from
	//so that eventually we are able to
	//look up info about the sensor
	// in the databases
	Name string

	//Value that the sensor has captured
	Value float64

	//there may be lags in receiving the
	//messages, this helps understand
	//when the reading was actually taken
	Timestamp time.Time
}

//can't send a raw type directly into rabbit.
//need to encode the data in some way
//Could have used json or protocol buffers for they are
//language neutral formats, but insted using gob (more efficient)
//as the application is pure GO.
func init() {
	//Registering gob package so it knows how
	//to work when we call on it later.

	//Every consumer of the package, can now rely on the SensorMessage
	//object being ready to send over the wire.
	gob.Register(SensorMessage{})
}
