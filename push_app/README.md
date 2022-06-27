Secrecy Kafka push client

This is the client application that lives on the OpenShift servers of RedHat clients.

Its duty is to periodically collect, preprocess and send logs to secrecy servers for further calculation.

Its development relies on Golang v1.18.2, confluent-kafka-go v1.8.2 and go-yaml v3

To install confluent kafka package, either run get_confluent_kafka_go.sh or run *go get github.com/confluentinc/confluent-kafka-go/kafka* directly in console.

- Assumptions made:
1. All kafka producers have the same timeouts (Can be made different in future)
2. All producers push to the same single topic (Maybe different multiple topics?)