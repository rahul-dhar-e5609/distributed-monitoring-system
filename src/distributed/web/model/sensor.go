package model

type Sensor struct {
	Name         string  `json:"name"`
	SerialNo     string  `json:"serialNo"`
	UnitType     string  `json:"unitType"`
	MinSafeValue float64 `json:"minSafeValue"`
	MaxSafeValue float64 `json:"maxSafeValue"`
}

func GetSensorByName(name string) (Sensor, error) {
	q := `SELECT name, serial_no, unit_type,
          min_safe_value, max_safe_value
        FROM sensor
        WHERE name = $1`

	result := Sensor{}

	row := db.QueryRow(q, name)
	err := row.Scan(&result.Name, &result.SerialNo, &result.UnitType,
		&result.MinSafeValue, &result.MaxSafeValue)

	return result, err
}
