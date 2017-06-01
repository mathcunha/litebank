package main

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

type KafkaProducer struct {
	c        chan *Event
	producer sarama.SyncProducer
}

type KafkaConsumer struct {
	c        chan *Event
	consumer sarama.PartitionConsumer
}

var (
	brokers = []string{"127.0.0.1:9092"}
	topic   = "litebank-transactions"
	topics  = []string{topic}
)

func init() {
	if kafka := os.Getenv("KAFKA_PORT"); kafka != "" {
		brokers[0] = kafka
		log.Printf("INFO: Kafka broker at %v \n", kafka)
	}
}

func newKafkaConfiguration() *sarama.Config {
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Return.Successes = true
	conf.ChannelBufferSize = 1
	conf.Version = sarama.V0_10_1_0
	return conf
}

func newKafkaSyncProducer() (sarama.SyncProducer, error) {
	kafka, err := sarama.NewSyncProducer(brokers, newKafkaConfiguration())

	if err != nil {
		return nil, err
	}

	return kafka, nil
}

func newKafkaConsumer() (sarama.Consumer, error) {
	consumer, err := sarama.NewConsumer(brokers, newKafkaConfiguration())

	if err != nil {
		return nil, err
	}

	return consumer, nil
}

func (p *KafkaProducer) start() {
	go func() {
		for {
			event := <-p.c
			sendEvent(p.producer, *event)
		}
	}()
}

func (p *KafkaConsumer) start() {
	go func() {
		for {
			e := consumeEvents(p.consumer)
			p.c <- e
		}
	}()
}

func consumeEvents(consumer sarama.PartitionConsumer) *Event {
	event := Event{}

	select {
	case err := <-consumer.Errors():
		log.Printf("Kafka error: %s\n", err)
	case msg := <-consumer.Messages():
		if err := json.Unmarshal(msg.Value, &event); err != nil {
			log.Printf("Failed parsing: %s", err)
		}
		return &event
	}
	return nil
}

func sendEvent(kafka sarama.SyncProducer, event interface{}) error {
	json, err := json.Marshal(event)

	if err != nil {
		log.Printf("Encoding error: %s\n", err)
		return err
	}

	msgLog := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(string(json)),
	}

	partition, offset, err := kafka.SendMessage(msgLog)
	if err != nil {
		log.Printf("Kafka error: %s\n", err)
		return err
	}

	log.Printf("Kafka - Message is stored in partition %d, offset %d\n", partition, offset)

	return nil
}
