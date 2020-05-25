---
description: messaging systems
title: messaging systems
---

### messaging systems

Push data into a queue, push (or pull) data out to clients.

- message log: store messages for future clients
- message queue: message delivered to 1 client
- pubsub: message delivered to all clients

#### _servers_

- Use GCP PubSub in the cloud.
- Use NATS if you don't need persistence.
- Cry and try Pulsar and/or NATS + Liftbridge/Jetstream if you need persistence

| queue                   | Language | Log/Queue/Pubsub | Delivery                    | Retention | Protocols                            | Operation modes | License    |
| ----------------------- | -------- | ---------------- | --------------------------- | --------- | ------------------------------------ | --------------- | ---------- |
| [ActiveMQ][activemq]    | Java     | Queue            | at most 1                   | no        | AMQP, HornetQ, MQTT, OpenWire, STOMP | cluster         | Apache 2.0 |
| [Kafka][kafka]          | Java     | Log/Queue        | at most/least/exactly 1     | yes       | Kafka                                | cluster         | Apache 2.0 |
| [NATS][nats]            | Go       | Pubsub           | at most 1                   | no        | NATS                                 | single          | Apache 2.0 |
| [NATS Streaming][natss] | Go       | Queue            | at least 1                  | yes       | NATS Streaming                       | cluster         | Apache 2.0 |
| [NSQ][nsq]              | Go       | Pubsub           | at least 1                  | no        | NSQ                                  | cluster         | MIT        |
| _[PubSub][pubsub]_      | -        | Pubsub           | at least 1                  | no        | HTTP, GRPC                           | hosted          | -          |
| [Pulsar][pulsar]        | Java     | Log/Queue/Pubsub | at most/least/effectively 1 | yes       | Pulsar                               | cluster         | Apache 2.0 |
| [RabbitMQ][rabbitmq]    | Erlang   | Queue            | at most/least 1             | yes       | AMQP, HTTP,MQTT,STOMP                | cluster         | MPL 1.1    |
| [RocketMQ][rocketmq]    | Java     | Queue            | at least 1                  | no        | RocketMQ, JMS, OpenMessaging         | cluster         | Apache 2.0 |

[activemq]: https://github.com/apache/activemq
[kafka]: https://github.com/apache/kafka
[nats]: https://github.com/nats-io/nats-server
[natss]: https://github.com/nats-io/nats-streaming-server
[nsq]: https://github.com/nsqio/nsq
[pubsub]: https://cloud.google.com/pubsub
[pulsar]: https://github.com/apache/pulsar
[rabbitmq]: https://github.com/rabbitmq/rabbitmq-server
[rocketmq]: https://github.com/apache/rocketmq

##### _notes_

- ActiveMQ optimizes for ???
- Kafka optimizes for throughput (streaming)
- Kafka is the recommended open source messaging system but it's hard to operate
- NATS optimizes for throughput and latency
- NATS is ephemeral
- NATS Streaming adds another layer over NATS for reliability, meh effectiveness
- NATS Streaming cluster = failover
- NSQ has weird delivery order
- NSQ optimizes for async
- NSQ cluster = failover with loss, or Raft with poor performance
- Liftbridge may be better alternative to NATS Streaming build on NATS
- Pulsar optimizes for latency
- Pulsar approx. next gen Kafka
- Pulsar is recommended over kafka but still hard to operate
- RabbitMQ has bad HA
- RabbitMQ is not recommended over single node
- RocketMQ approx. next gen ActiveMQ
- RocketMQ optimizes for batch processing
- RocketMQ has poor docs

#### _other_

[ZeroMQ][zeromq] is a set of protocol/libraries for direct messaging between producers/consumers.

[zeromq]: https://github.com/zeromq
