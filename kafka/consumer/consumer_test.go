package consumer_test

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/labstack/gommon/log"
	"github.com/stretchr/testify/assert"
	"synergetic-craft/kafka/consumer"
	"synergetic-craft/kafka/producer"
	"testing"
	"time"
)

const (
	ErrConnectFixture = "failed connect"
)

func TestConsumer_Connect(t *testing.T) {
	t.Parallel()

	t.Run("should return success when new consumer is created", func(t *testing.T) {
		mockProducer, _ := kafka.NewMockCluster(1)
		defer mockProducer.Close()

		broker := mockProducer.BootstrapServers()

		c := consumer.NewConsumer(broker, "groupTest", "test", true)
		err := c.Connect()

		assert.NoError(t, err)
	})

	t.Run("should return error when consumer failed create", func(t *testing.T) {
		mockProducer, _ := kafka.NewMockCluster(1)
		defer mockProducer.Close()

		broker := mockProducer.BootstrapServers()

		c := consumer.NewConsumer(broker, "groupTest", "", true)
		err := c.Connect()

		assert.Error(t, err)
		assert.Equal(t, "Local: Invalid argument or configuration", err.Error())
	})
}

func TestConsumer_EventProcessor(t *testing.T) {
	t.Parallel()

	t.Run("should read message and resolve with the assigned function", func(t *testing.T) {
		mockProducer, _ := kafka.NewMockCluster(1)
		defer mockProducer.Close()

		broker := mockProducer.BootstrapServers()

		c := consumer.NewConsumer(broker, "group", "test", true)
		err := c.Connect()
		if err != nil {
			t.Fatal(ErrConnectFixture)
		}

		p, _ := producer.NewProducer(broker, 2000)
		_ = p.Connect()

		_ = p.Send("test", []byte(`key`), []byte(`{"name":"new event","param":"hi, how are you?"}`))

		c.SetHandlers(map[string]func([]byte) error{
			"new event": func(event []byte) error {
				var response struct {
					Name  string `json:"name"`
					Param string `json:"param"`
				}
				if err := json.Unmarshal(event, &response); err != nil {
					return err
				}

				log.Info(response.Param)

				return nil
			},
		})

		go func() {
			c.EventProcessor()
		}()

		time.Sleep(time.Duration(6) * time.Second)

		c.Stop()
	})
}
