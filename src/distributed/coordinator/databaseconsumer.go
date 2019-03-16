package coordinator

import (
	"distributed/dto"
	"distributed/qutils"

	"github.com/streadway/amqp"

	"bytes"
	"encoding/gob"
	"time"
)

const maxRate = 5 * time.Second

// DatabaseConsumer is a struct which is
// responsible for listening to events
// that are being emitted throughtout the
// system and decide which ones to forward
// to the database manager
type DatabaseConsumer struct {
	er EventRaiser //  we can pass callback without this type knowing about how to publish events itself

	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue // Queue we route messages to that tht DataConsumer decides should be saved

	sources []string
}

// NewDatabaseConsumer is a constructor function
// for DatabaseConsumer, using EventRaiser to
// add listners (Dependency Injection)
func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {

	dc := DatabaseConsumer{
		er: er,
	}

	dc.conn, dc.ch = qutils.GetChannel(url)

	dc.queue = qutils.GetQueue(qutils.PersistReadingsQueue,
		dc.ch, false)

	// Added listener for the event that gets published
	// whenever a new data spurce is discovered
	dc.er.AddListener("DataSourceDiscoveredEvents", func(eventData interface{}) {
		// Subscribing to the newly discovered
		// data source
		dc.SubscribeToDataEvent(eventData.(string))
	})

	return &dc
}

// SubscribeToDataEvent func is reponsible for subscribing
// to the events published.
func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {

	// checking if event has been subscriber to already
	for _, v := range dc.sources {
		//bail, if already subscribed
		if v == eventName {
			return
		}
	}

	// Subscribing new data source
	// and listening for messages received
	// event.
	//
	// Using closure for returning callback
	// Gives a new isolated variable scope
	// that gets created every time the func
	// is called
	//
	// This allows the event handlers to register
	// a state that is needed to throttle down the
	// rate with which the messages are coming in
	// from the event sources.
	dc.er.AddListener("MessageReceived_"+eventName, func() func(interface{}) {

		// state for persisting data
		//
		// retains their state from call to call
		// since theire state is captured by the closure
		// not the callback
		prevTime := time.Unix(0, 0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {

			// casting event data object in
			// EventData struct
			ed := eventData.(EventData)

			// throttle the messages to be persisted
			// no faster than 5 seconds (maxrate)
			if time.Since(prevTime) > maxRate {

				// resetting the clock
				prevTime = time.Now()

				sm := dto.SensorMessage{
					Name:      ed.Name,
					Value:     ed.Value,
					Timestamp: ed.Timestamp,
				}

				// resetting buffer to make sure there is
				// no data from a previous call
				//
				// resetting state
				buf.Reset()

				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{
					Body: buf.Bytes(),
				}

				// Publishing to persist reading queue
				// to persist data in the database.
				dc.ch.Publish(
					"",
					qutils.PersistReadingsQueue,
					false,
					false,
					msg)
			}
		}
	}())
}
