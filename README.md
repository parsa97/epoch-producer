# Epoch producer

This project start producing time in epoch to kafka topic.

This app contains two part:
- Producer: Produce time in epch(ms) format to kafka topic.
- Metrics: collect metrics from producer and consumer.

### Config Global

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `LOG_LEVEL` | log level | enum(trace, debug, info, warn, error, fatal, panic) | info

### Config Exporter

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `METRICS_PATH` | metrics path | string | /metrics
| `LISTEN_PORT` | exporter listening port  | number | 8080


### Config Producer

| Name                                 | Description                                               |  Type | Default
|:-------------------------------------|:----------------------------------------------------------|:-----:|:--------:|
| `PRODUCER_VERSION` | Kafka Version | string | 2.7.0
| `TOPIC` | Kafka Topic to Produce | string | output
| `PRODUCER_BROKERS` | Kafka Brokers | comma seperated hosts:port | 127.0.0.1:9092
| `PRODUCER_MAX_MESSAGE_BYTES` | Producer.MaxMessageBytes | number | 1000000
| `PRODUCER_FLUSH_FREQUENCY` | Producer.Flush.Frequency | number(milisecound) | -
| `PRODUCER_FLUSH_MESSAGE` | Producer.Flush.Messages | number | -
| `PRODUCER_FLUSH_MAX_MESSAGE` | Producer.Flush.MaxMessages | number | -
| `PRODUCER_RETURN_SUCCESS` | Producer.Return.Successes | bool | true
| `PRODUCER_TIMEOUT` | Producer.Timeout | number(secound) | 10
| `PRODUCER_RETRY_MAX` | Producer.Retry.Max | number | 3
| `PRODUCER_RETRY_BACKOFF` | Producer.Retry.Backoff | number(milisecound) | 100
| `PRODUCER_RETURN_ERROR` | Producer.Return.Errors | bool | true
| `PRODUCER_COMPRESSIONLEVEL` | Producer.CompressionLevel | enum(gzip, zstd, snappy, lz4, none) | none
| `PRODUCER_PARTITIONER` | Producer.Partitioner | enum(random, hash, rr) | hash
| `PRODUCER_REQUIRED_ACKS` | Producer.RequiredAcks | enum(0, 1, -1) | 1
| `PRODUCER_CLIENTID` | Producer ClientID  | string | defaultClientID
| `PRODUCER_CHANNELBUFFERSIZE` | ChannelBufferSize | number | 256

**Note:** Required acks enums map to:
- 0: NoResponse
- 1: WaitForLocal
- -1: WaitForAll

Almost all descriptions are map to sarama
[config.go](https://github.com/Shopify/sarama/blob/master/config.go#L441)
values


# Metrics

The following metrics are available:

|Name|Description|
|---|---|
|`produced_message_total`|How many epoch message produced|

Metrics are counters and might be used with
[`rate()`](https://prometheus.io/docs/prometheus/latest/querying/functions/#rate())
to calculate per-second average rate.

Tests
-----
Run tests are require a reachable kafka broker
```
go test ./...
```

Build
-----
```
docker build .
[...]
Successfully built b334267f4035
```

Run container
-------------
```
docker run -p 8080:8080 -e KAFKA_BROKER=localhost:9092 b334267f4035
```

Test
----
Run Kafka console consumer on topic ouput
```
kafka-console-consumer --bootstrap-server localhost:9092 --topic input
```
verify epoch Messages
```
1615035196527
1615035196529
1615035196531
[...]
```