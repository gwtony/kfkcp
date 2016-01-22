package kfkp

import (
	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	addr		string
	log			*Log
	producer	sarama.SyncProducer
	consumer	sarama.Consumer

	/* TODO: should use interface */
	s			*Server
}

func InitKafka(addr string, flag bool, log *Log) (*KafkaClient, error) {
	kfk := &KafkaClient{}

	kfk.addr = addr
	kfk.log = log

	if flag {
		err := kfk.initKafkaProducer(addr)
		 if err != nil {
			 kfk.log.Error("Init kafka producer failed")
			 return nil, err
		 }
	} else {
		err := kfk.initKafkaConsumer(addr)
		 if err != nil {
			 kfk.log.Error("Init kafka consumer failed")
			 return nil, err
		 }
	}
	//TODO: producer & consumer to be closed

	return kfk, nil
}

func (k *KafkaClient) initKafkaProducer(broker string) error {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
    config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message

	brokerList := []string {broker}
	k.log.Debug("Call new sync producer: ")

    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        k.log.Error("Failed to start producer:", err)
		return err
    }

    k.producer = producer

	return nil
}

func (k *KafkaClient) initKafkaConsumer(broker string) error {
    config := sarama.NewConfig()
	config.Consumer.Return.Errors = true

	brokerList := []string {broker}
	k.log.Debug(brokerList)
    consumer, err := sarama.NewConsumer(brokerList, config)
    if err != nil {
        k.log.Error("Failed to start consumer:")
		k.log.Error(err)
		return err
    }

	k.log.Debug("Create new consumer done")
    k.consumer = consumer

	return nil
}

func (k *KafkaClient) SendData(topic string, msg string) error {
	k.log.Debug("SendData topic: %s, msg: %s", topic, msg)

	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})

	k.log.Debug("Send data to kafka done")
	if err != nil {
		k.log.Error("Kafka client send message failed", err)
		return err
	}

	return nil
}

