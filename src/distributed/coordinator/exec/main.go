package main

import (
	"distributed/coordinator"
	"fmt"
)

var dc *coordinator.DatabaseConsumer

func main() {

	// wiring event aggreagtor to coordinator package
	ea := coordinator.NewEventAggregator()
	// instantiating package level databse consumer
	dc = coordinator.NewDatabaseConsumer(ea)

	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
