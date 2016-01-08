package kfkp

import (
	"github.com/Shopify/sarama"
)

type KafkaClient struct {
	addr		string
	port		string
	log			*Log
	producer	sarama.SyncProducer
}

func InitKafka(addr string, port string, log *Log) (*KafkaClient, error) {
	kfk := &KafkaClient{}
	kfk.addr = addr
	kfk.port = port
	kfk.log = log
	err := kfk.initKafkaProducer(addr + ":" + port)
	if err != nil {
		kfk.log.Error("Init kafka producer failed")
		return nil, err
	}
	//kconsumer := kfk.initKafkaConsumer(addr + ":" + port)

	return kfk, nil
}

func (k *KafkaClient) initKafkaProducer(broker string) error {
    config := sarama.NewConfig()
    config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
    config.Producer.Retry.Max = 3                    // Retry up to 10 times to produce the message

	brokerList := []string {broker}
	k.log.Debug("call new sync producer: ")
	k.log.Debug(brokerList, config)

    producer, err := sarama.NewSyncProducer(brokerList, config)
    if err != nil {
        k.log.Error("Failed to start producer:", err)
		return err
    }

    k.producer = producer

	return nil
}

func (k *KafkaClient) sendData(topic string, msg string) error {
	k.log.Debug("[sendData] topic: %s, msg: %s", topic, msg)

	//partition, offset, err := k.producer.SendMessage(&sarama.ProducerMessage{
	_, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	})

	k.log.Debug("send data to kafka done")
	if err != nil {
		k.log.Error("Kafka client send message failed", err)
		return err
	}

	return nil
}
