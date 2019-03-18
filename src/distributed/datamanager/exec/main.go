package main

import (
	"bytes"
	"distributed/datamanager"
	"distributed/dto"
	"distributed/qutils"
	"encoding/gob"
	"log"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	conn, ch := qutils.GetChannel(url)
	defer conn.Close()
	defer ch.Close()

	// receiving messages from the persist reading queue
	msgs, err := ch.Consume(
		qutils.PersistReadingsQueue, // queue
		"",                          // consumer, no need for a new name, letting rabbit to assign a new name
		false,                       // autoack, till then this has been set to true, so that the messages are acknowledged and removed  from the queue as soon as they
		// are consumed, but this time, we want to take an advantage of the fact that rabbit has these messages in memory to build a simple transactional system for the db,
		// by setting this to false, we now have to manually acknowledge the message, which that i can wait up until i have verified that the databse has successfully saved a record,
		// if somehting goes wrong then we simply dont ack the message and rabbit wikk try and deliver it again later.
		true, // exclusive, every module in this app is designed to support multiple instances running at the same time, while we could probably do that
		// for the datamanager as well, by setting this to true, it is gonna fail to connect if it cant get an exclusive connection to this queue. Practiaclly meaning we can only have one database manager running at a simgle time.
		false, //nolocal
		false, // nowait
		nil)   //args amqp.Table
	if err != nil {
		log.Fatal("Failed to get access to messages")
	}

	for msg := range msgs {
		// Decoding the message
		// - wrapping the message body with a buffer
		// - use the buffer to create a decoder
		// - placeholder sensor message object that holds the decoded reading
		// - using decoder's decode method to populate it
		buf := bytes.NewReader(msg.Body)
		dec := gob.NewDecoder(buf)
		sd := &dto.SensorMessage{}
		dec.Decode(sd)

		err := datamanager.SaveReader(sd)

		if err != nil {
			log.Printf("Failed to save reading from sensor %v. Error: %s",
				sd.Name, err.Error())
		} else {
			// acknowledge that message was processed properly
			msg.Ack(false)
		}
	}
}
