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
		entity, err := event.loadEntity()
		if err != nil {
			return err
		}
		if event.Type == NewTransactionEvent {
			ok, err := processTransaction(entity.(*Transaction))
			if err != nil {
				return err
			}
			log.Printf("{\"EventType\":%b - %T, %v}\n", event.Type, entity, ok)
			continue
		}
		if (event.Type & New) == New {
			err = insert(entity.collection(), entity)
			if err != nil {
				return err
			}
			log.Printf("{\"EventType\":%b - %T}\n", event.Type, entity)

		} else {

			log.Println("not persisted ", event)
		}
	}

	return nil
}
