package datamanager

import (
	"distributed/dto"
	"errors"
)

// To cache the relationship between the sensor's
// names and theire IDs in the database so that we
// need not constantly query the database
var sensors map[string]int

// SaveReader is a save function that accepts a pointer
// to a sensor message object and return an error object
func SaveReader(reading *dto.SensorMessage) error {
	// checking if already have the ID for the sensor
	if sensors[reading.Name] == 0 {
		// not found in cache, fetching relationship
		getSensors()
	}

	// checking if current sensor's name is mapped
	// after querying db
	if sensors[reading.Name] == 0 {
		return errors.New("Unable to find the sensor for name '" +
			reading.Name + "'")
	}

	// insert reading, of sensor exists
	q := `
		INSERT INTO sensor_reading
		(value, senosr_id, taken_on)
		VALUES
		($1, $2, $3)
	`

	_, err := db.Exec(q, reading.Value, sensors[reading.Name], reading.Timestamp)

	return err

}

func getSensors() {
	// Re- initializing the map so that
	// any stale data is discarded
	sensors = make(map[string]int)

	q := `SELECT id, name FROM sensor`

	// Querying the database
	rows, _ := db.Query(q)

	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		// updating map
		sensors[name] = id
	}
}
