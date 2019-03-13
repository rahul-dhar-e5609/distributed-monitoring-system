package coordinator

import "time"

// EventAggregator struct
type EventAggregator struct {
	// Map of consumers
	//
	// Lets anything to register as an event listener
	// That is done by providing the event it is interested
	// in and the callback function that handles the event
	listeners map[string][]func(EventData)
}

// NewEventAggregator is the constructor
// function for the EventAggregator
func NewEventAggregator() *EventAggregator {
	ea := EventAggregator{
		listeners: make(map[string][]func(EventData)),
	}

	return &ea
}

// AddListener maps a callback function to the event
//
// @param name | string - name of the event being registered for
// @param f | func - callback
func (ea *EventAggregator) AddListener(name string, f func(EventData)) {
	ea.listeners[name] = append(ea.listeners[name], f)
}

// PublishEvent function is responsible for publishing
// the event to all the listeners
//
// @param name | string - name of the event
// @param eventData | EventData - data
func (ea *EventAggregator) PublishEvent(name string, eventData EventData) {
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
