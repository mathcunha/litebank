package main

import (
	"github.com/Shopify/sarama"
	"log"
)

func consume() error {
	kafka, err := newKafkaConsumer()
	if err != nil {
		return err
	}
	defer kafka.Close()

	consumer, err := kafka.ConsumePartition(topic, 0, sarama.OffsetOldest)
	if err != nil {
		return err
	}
	kConsumer = KafkaConsumer{make(chan *Event, 10), consumer}
	kConsumer.start()

	for {
		event := <-kConsumer.c
		log.Println(event)
	}

	return nil
}
