package logmill

import (
	"os"
	"fmt"
	"time"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaConf configuration 
type KafkaConf struct {
	Broker string
	Topic string
}

// NewKafkaConf return the default kafka configuration
func NewKafkaConf() *KafkaConf {

	return &KafkaConf {
		Broker :  "localhost:9092",
		Topic : "test",
	}

}

// KafkaLogmill sends logs to the kafka server
type KafkaLogmill struct {
	Mill
	KafkaConf
	kafkaProducer *kafka.Producer
	deliveryChan  chan kafka.Event
}

// NewKafkaLogmill creates a new Kafka Logmill
func NewKafkaLogmill( kafkaConf *KafkaConf, lg LogGenerator) *KafkaLogmill {

	kafkaLogmill := new(KafkaLogmill)

	kafkaLogmill.lg = lg

	kafkaLogmill.Broker = kafkaConf.Broker
	kafkaLogmill.Topic = kafkaConf.Topic
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": kafkaConf.Broker})
	kafkaLogmill.kafkaProducer = p

	dc := make(chan kafka.Event)
	kafkaLogmill.deliveryChan = dc

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}

	return kafkaLogmill
}

// SendLogs send the logs per the configuration
func (f *KafkaLogmill) SendLogs( tick time.Duration, logsPerTick int , nLogsToSend int ) {

	sendClock(f, tick, logsPerTick, nLogsToSend)

	close(f.deliveryChan)

}

func (f *KafkaLogmill) logGenerator() LogGenerator {
	return f.lg
}

func (f *KafkaLogmill) writeLog(logText string) (bytesSent int64) {


	err := f.kafkaProducer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &f.Topic, Partition: kafka.PartitionAny},
		Value:          []byte(logText),
		Headers:        []kafka.Header{{Key: "myTestHeader", Value: []byte("header values are binary")}},
	}, f.deliveryChan)

	if err != nil {
		fmt.Printf("kafkaProducer.Produce err = %v\n",err)
	}
	e := <-f.deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	return int64(len(logText))
}
