package producer

import (
	"errors"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"time"
)

type Producer interface {
	Connect() error
	Close()
	Send(topic string, key, message []byte) <-chan error
}

type ParamConn struct {
	broker         string
	timeoutMessage time.Duration
}

type producer struct {
	paramConn ParamConn
	producer  *kafka.Producer
}

func NewProducer(broker string, timeoutMessage int64) (Producer, error) {
	if broker == "" || timeoutMessage == 0 {
		return nil, errors.New("parameters cannot be zero")
	}

	return &producer{
		paramConn: ParamConn{
			broker:         broker,
			timeoutMessage: time.Duration(timeoutMessage) * time.Millisecond,
		},
	}, nil
}

func (k *producer) Send(topic string, key, message []byte) <-chan error {
	responseChan := make(chan error)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          message,
		Key:            key,
	}

	go func() {
		deliveryChan := make(chan kafka.Event)

		err := k.producer.Produce(msg, deliveryChan)
		if err != nil {
			responseChan <- fmt.Errorf("error producing message: %v", err)
			close(responseChan)

			return
		}

		select {
		case <-time.After(k.paramConn.timeoutMessage):
			responseChan <- errors.New("kafka message [time exceeded]")
			close(responseChan)

		case result := <-deliveryChan:
			msgResponse := result.(*kafka.Message)

			if msgResponse.TopicPartition.Error != nil {
				responseChan <- msgResponse.TopicPartition.Error
			} else {
				responseChan <- nil
			}

			close(responseChan)

			log.Infof("kafka message [topic = %s] [value = %s]", msgResponse.TopicPartition.Topic, msgResponse.Value)
		}
	}()

	return responseChan
}

func (k *producer) Connect() error {
	conn, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": k.paramConn.broker,
	})
	if err != nil {
		return err
	}

	k.producer = conn

	return nil
}

func (k *producer) Close() {
	k.producer.Close()
}
