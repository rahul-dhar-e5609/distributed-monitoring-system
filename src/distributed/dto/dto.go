package dto

import "time"

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
