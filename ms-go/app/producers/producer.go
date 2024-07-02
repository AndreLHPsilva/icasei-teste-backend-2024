package producers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

type ProductWithoutTimestamps struct {
	ID          int       `json:"id" bson:"id,omitempty"`
	Name        string    `json:"name" bson:"name,omitempty"`
	Brand       string    `json:"brand"  bson:"brand,omitempty"`
	Price       float64   `json:"price"  bson:"price,omitempty"`
	Description string    `json:"description"  bson:"description,omitempty"`
	Stock       int       `json:"stock"  bson:"stock,omitempty"`
}

type KafkaPayload struct {
	Product ProductWithoutTimestamps `json:"product"`
}

func ProducerKafka(payload KafkaPayload) error {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  []string{"kafka:29092"},
		Topic:    "go-to-rails",
	})

	defer writer.Close()

	message, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:     nil,
		Value:   message,
		Headers: nil,
	})
	if err != nil {
		return err
	}

	log.Println("Mensagem enviada com sucesso para o Kafka")
	return nil
}
