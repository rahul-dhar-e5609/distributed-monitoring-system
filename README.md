# Distributed Monitoring System

> A boilerplate to checkout distributed programming in Go.

The queue discovery was initially a direct exchange. The sensors published two kinds of messages.

- **Data message** - Contains the name, timestamp and the value of the reading that it took. This data is published onto a route that is created by the sensor at runtime and is not discovered before that. To address that there is another message that is published onto a known route and publishes the routing key of the data queue.

- **Route message** - It messages the routing key onto a known route. Since this route is well known, it becomes easier for the downstream coordinators to know about the data routes.

The data routes are published to the default exchange, which is a direct exchange, that means only one consumer receives the message. This ensures no matter how many coordinators are active, the message broker guarantees that only one consumer gets the message at only one time.

### TODO

- **Discovery route**: All the coordinators should know when a sensor comes online, direct exchange fails this requirement, instead a fanout exchange can inform all the attached queues about the received message, doing so, all the coordinators can know about the received message.
- **Fanout Exchange for lazy coordinators**: Need another fanout exchange for the coordinators that spawned after all the sensors got instantiated. The coordinators that spawned late have no way to know about the sensors that got instantiated, there should be a fanout exchange that flows reverse to other data flows. The coordinators can then make a discovery request to the exchange, that message then fans out to all the sensors which then respond by publishing there data queue's name to the fanout exchange.
- **Event Bus**: Modify the Event Aggregator to use channels instead of event pattern. 
Check out Simple Event Bus in go by Kasun Vithanage
https://levelup.gitconnected.com/lets-write-a-simple-event-bus-in-go-79b9480d8997
