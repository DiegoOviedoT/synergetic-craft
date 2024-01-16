package consumer

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
)

type Consumer interface {
	Stop()
	Connect() error
	EventProcessor()
	SetHandlers(handlers map[string]func([]byte) error)
}

type consumer struct {
	topic         string
	groupID       string
	broker        string
	handlers      map[string]func([]byte) error
	enableLogging bool
	isRunning     bool
	consumer      *kafka.Consumer
}

type message struct {
	Name string `json:"name"`
}

func NewConsumer(broker string, groupID, topic string, enableLogging bool) Consumer {
	return &consumer{
		topic:         topic,
		groupID:       groupID,
		broker:        broker,
		enableLogging: enableLogging,
		isRunning:     true,
	}
}

func (kc *consumer) Connect() error {
	config := kafka.ConfigMap{
		"bootstrap.servers": kc.broker,
		"group.id":          kc.groupID,
		"auto.offset.reset": "earliest",
	}

	var err error

	kc.consumer, err = kafka.NewConsumer(&config)
	if err != nil {
		return err
	}

	if err = kc.consumer.SubscribeTopics([]string{kc.topic}, nil); err != nil {
		kc.Stop()

		return err
	}

	return nil
}

func (kc *consumer) EventProcessor() {
	defer kc.Stop()

	for kc.isRunning {
		event := kc.consumer.Poll(100)
		if event == nil {
			continue
		}

		err := kc.event(event)

		if kc.enableLogging {
			log.Info(err)
		}
	}
}

func (kc *consumer) SetHandlers(handlers map[string]func([]byte) error) {
	kc.handlers = handlers
}

func (kc *consumer) Stop() {
	if kc.consumer != nil {
		_ = kc.consumer.Close()
	}

	kc.isRunning = false
}

func (kc *consumer) event(event kafka.Event) (errEvent error) {
	switch ev := event.(type) {
	case *kafka.Message:
		var msg message

		if err := json.Unmarshal(ev.Value, &msg); err != nil {
			return err
		}

		if _, ok := kc.handlers[msg.Name]; !ok {
			return fmt.Errorf("event handler [ %s ] does not exist", msg.Name)
		}

		return kc.handlers[msg.Name](ev.Value)
	case kafka.Error:
		return fmt.Errorf("error code [ %v ]\nevent [ %v ]", ev.Code(), ev)
	}

	return nil
}
