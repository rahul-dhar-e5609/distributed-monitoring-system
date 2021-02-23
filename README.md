# Distributed Monitoring System

[![Go Report Card](https://goreportcard.com/badge/github.com/IAmRDhar/distributed-monitoring-system)](https://goreportcard.com/report/github.com/IAmRDhar/distributed-monitoring-system)

> A distributed system in Go using RabbitMQ Exchanges and Queues for monitoring the communications among data generators, data coordinators and data consumers.

The queue discovery was initially a direct exchange. The sensors published two kinds of messages.

- **Data message** - Contains the name, timestamp and the value of the reading that it took. This data is published onto a route that is created by the sensor at runtime and is not discovered before that. To address that there is another message that is published onto a known route and publishes the routing key of the data queue.

- **Route message** - It messages the routing key onto a known route. Since this route is well known, it becomes easier for the downstream coordinators to know about the data routes.

The data routes are published to the default exchange, which is a direct exchange, that means only one consumer receives the message. This ensures no matter how many coordinators are active, the message broker guarantees that only one consumer gets the message at only one time.
