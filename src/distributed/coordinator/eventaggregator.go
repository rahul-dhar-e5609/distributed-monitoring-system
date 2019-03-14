package coordinator

import "time"

// EventRaiser interface is used to
// prevent event consumers from knowing
// how to publish events themselves.
//
// Exposing as little functionality as necessery
// using the interface that is implemented by the
// event aggregator
type EventRaiser interface {
	// changing the data that the callback receives
	// from an EventData object to empty interface
	// so that there is loose coupling and any kind of data
	// can be sent.
	AddListener(eventName string, f func(interface{}))
}

// EventAggregator struct
type EventAggregator struct {
	// Map of consumers
	//
	// Lets anything to register as an event listener
	// That is done by providing the event it is interested
	// in and the callback function that handles the event
	listeners map[string][]func(interface{})
}

// NewEventAggregator is the constructor
// function for the EventAggregator
func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(interface{})),
	}

	return &ea
}

// AddListener maps a callback function to the event
//
// @param name | string - name of the event being registered for
// @param f | func - callback
func (ea *EventAggregator) AddListener(name string, f func(interface{})) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

// PublishEvent function is responsible for publishing
// the event to all the listeners
//
// @param name | string - name of the event
// @param eventData | EventData - data
func (ea *EventAggregator) PublishEvent(name string, eventData interface{}) {
	if ea.listeners[name] != nil {
		for _, r := range ea.listeners[name] {
			// not sending a pointer to eventData
			// therefore sends a copy to all the
			// consumers and not the actual object.
			// Consumers dont fiddling the data and
			// confuse one another.
			r(eventData)
		}
	}
}

//EventData struct
type EventData struct {
	Name      string
	Value     float64
	Timestamp time.Time
}
