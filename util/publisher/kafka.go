package publisher

import (
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
)

var kafkaPublisher *kafka.Publisher

func GetKafkaPublisher() *kafka.Publisher {
	return kafkaPublisher
}

func InitKafkaPublisher(brokers []string) error {
	watermillLogger := watermill.NewStdLogger(true, true)

	// Publisher
	publisher, err := kafka.NewPublisher(kafka.PublisherConfig{Brokers: brokers, Marshaler: kafka.DefaultMarshaler{}}, watermillLogger)
	if err != nil {
		return err
	}

	kafkaPublisher = publisher

	return nil
}
