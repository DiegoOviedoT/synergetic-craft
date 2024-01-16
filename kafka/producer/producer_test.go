package producer_test

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	kafkaLocal "synergetic-craft/kafka/producer"
	"testing"
)

const (
	broker            = "localhost:9093"
	ErrConnectFixture = "failed connect"
)

func TestNewProducer(t *testing.T) {
	t.Run("should return new producer when parameters is send", func(t *testing.T) {
		p, err := kafkaLocal.NewProducer(broker, 2000)

		assert.NoError(t, err)
		assert.NotNil(t, p)
	})

	t.Run("should return error when parameters are not send", func(t *testing.T) {
		p, err := kafkaLocal.NewProducer("", 0)

		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestProducer_Connect(t *testing.T) {
	t.Run("should return connection when producer is created correctly", func(t *testing.T) {
		p, _ := kafkaLocal.NewProducer(broker, 2000)
		defer p.Close()

		err := p.Connect()

		assert.NotNil(t, p)
		assert.NoError(t, err)
	})
}

func TestProducer_Send(t *testing.T) {
	t.Run("should return success when message is send", func(t *testing.T) {
		mockProducer, _ := kafka.NewMockCluster(1)
		defer mockProducer.Close()

		broker := mockProducer.BootstrapServers()

		p, _ := kafkaLocal.NewProducer(broker, 2000)
		err := p.Connect()
		if err != nil {
			t.Fatal(ErrConnectFixture)
		}

		errChan := p.Send("test", []byte(`key`), []byte(`{"name":"new event"}`))

		err = <-errChan

		assert.NoError(t, err)
	})

	t.Run("should return timeout when message exceeded timeout", func(t *testing.T) {
		p, _ := kafkaLocal.NewProducer(broker, 10)
		defer p.Close()

		err := p.Connect()
		if err != nil {
			t.Fatal(ErrConnectFixture)
		}

		errChan := p.Send("test", []byte(`key`), []byte(`{"name":"event time out"}`))

		err = <-errChan

		assert.Error(t, err)
		assert.Equal(t, "kafka message [time exceeded]", err.Error())
	})

	t.Run("should return err when topic is not defined", func(t *testing.T) {
		mockProducer, _ := kafka.NewMockCluster(1)
		defer mockProducer.Close()

		broker := mockProducer.BootstrapServers()

		p, _ := kafkaLocal.NewProducer(broker, 2000)
		err := p.Connect()
		if err != nil {
			t.Fatal(ErrConnectFixture)
		}

		errChan := p.Send("", []byte(`key`), []byte(`{"name":"new event"}`))

		err = <-errChan

		assert.Error(t, err)
		assert.Equal(t, "error producing message: Local: Invalid argument or configuration", err.Error())
	})
}
