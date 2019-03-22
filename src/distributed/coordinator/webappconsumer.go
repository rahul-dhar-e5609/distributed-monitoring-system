package coordinator

import (
	"bytes"
	"distributed/dto"
	"distributed/qutils"
	"encoding/gob"
	"github.com/streadway/amqp"
)

type WebappConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources []string
}

func NewWebappConsumer(er EventRaiser) *WebappConsumer {
	wc := WebappConsumer{
		er: er,
	}

	wc.conn, wc.ch = qutils.GetChannel(url)
	qutils.GetQueue(qutils.PersistReadingsQueue, wc.ch, false)

	go wc.ListenForDiscoveryRequests()

	wc.er.AddListener("DataSourceDiscovered",
		func(eventData interface{}) {
			wc.SubscribeToDataEvent(eventData.(string))
		})

	wc.ch.ExchangeDeclare(
		qutils.WebappSourceExchange, //name string,
		"fanout",                    //kind string,
		false,                       //durable bool,
		false,                       //autoDelete bool,
		false,                       //internal bool,
		false,                       //noWait bool,
		nil)                         //args amqp.Table)

	wc.ch.ExchangeDeclare(
		qutils.WebappReadingsExchange, //name string,
		"fanout",                      //kind string,
		false,                         //durable bool,
		false,                         //autoDelete bool,
		false,                         //internal bool,
		false,                         //noWait bool,
		nil)                           //args amqp.Table)

	return &wc
}

func (wc *WebappConsumer) ListenForDiscoveryRequests() {
	q := qutils.GetQueue(qutils.WebappDiscoveryQueue, wc.ch, false)
	msgs, _ := wc.ch.Consume(
		q.Name, //queue string,
		"",     //consumer string,
		true,   //autoAck bool,
		false,  //exclusive bool,
		false,  //noLocal bool,
		false,  //noWait bool,
		nil)    //args amqp.Table)

	for range msgs {
		for _, src := range wc.sources {
			wc.SendMessageSource(src)
		}
	}
}

func (wc *WebappConsumer) SendMessageSource(src string) {
	wc.ch.Publish(
		qutils.WebappSourceExchange, //exchange string,
		"",    //key string,
		false, //mandatory bool,
		false, //immediate bool,
		amqp.Publishing{Body: []byte(src)}) //msg amqp.Publishing)
}

func (wc *WebappConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range wc.sources {
		if v == eventName {
			return
		}
	}

	wc.sources = append(wc.sources, eventName)

	wc.SendMessageSource(eventName)

	wc.er.AddListener("MessageReceived_"+eventName,
		func(eventData interface{}) {
			ed := eventData.(EventData)
			sm := dto.SensorMessage{
				Name:      ed.Name,
				Value:     ed.Value,
				Timestamp: ed.Timestamp,
			}

			buf := new(bytes.Buffer)

			enc := gob.NewEncoder(buf)
			enc.Encode(sm)

			msg := amqp.Publishing{
				Body: buf.Bytes(),
			}

			wc.ch.Publish(
				qutils.WebappReadingsExchange, //exchange string,
				"",    //key string,
				false, //mandatory bool,
				false, //immediate bool,
				msg)   //msg amqp.Publishing)
		})
}
