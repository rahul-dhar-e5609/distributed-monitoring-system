package coordinator

import (
	"distributed/dto"
	"distributed/qutils"

	"github.com/streadway/amqp"

	"bytes"
	"encoding/gob"
	"time"
)

// Responsible for listening to events
// that are being emitted throughtout the
// system and decide which ones to forward
// to the database manager

const maxRate = 5 * time.Second

type DatabaseConsumer struct {
	er EventRaiser //  we can pass callback without this type knowing about how to publish events itself

	conn  *amqp.Connection
	ch    *amqp.Channel
	queue *amqp.Queue // Queue we route messages to that tht DataConsumer decides should be saved

	sources []string
}

// NewDataConsumer is a constructor function
// for DatabaseConsumer, using EventRaiser to
// add listners (Dependency Injection)
func NewDataConsumer(er EventRaiser) *DatabaseConsumer {

	dc := DatabaseConsumer{
		er: er,
	}

	dc.conn, dc.ch = qutils.GetChannel(url)

	dc.queue = qutils.GetQueue(qutils.PersistReadingsQueue,
		dc.ch, false)

	dc.er.AddListener("DataSourceDiscoveredEvents", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})

	return &dc
}

// SubscribeToDataEvent
func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range dc.sources {
		if v == eventName {
			return
		}
	}

	// registering new data source
	dc.er.AddListener("MessageReceived_"+eventName, func() func(interface{}) {
		prevTime := time.Unix(0, 0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			ed := eventData.(EventData)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()

				sm := dto.SensorMessage{
					Name:      ed.Name,
					Value:     ed.Value,
					Timestamp: ed.Timestamp,
				}

				buf.Reset()

				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{
					Body: buf.Bytes(),
				}

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
