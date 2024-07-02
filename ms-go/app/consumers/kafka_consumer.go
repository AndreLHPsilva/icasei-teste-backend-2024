package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"ms-go/app/models"
	"ms-go/app/services/products"
	"github.com/segmentio/kafka-go"
)

type KafkaPayload struct {
	Product   models.Product `json:"product"`
	Operation string         `json:"operation"`
}

func ConsumerKafka() {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"kafka:29092"},
		GroupID: "consumer",
		Topic:   "rails-to-go",
	})

	defer reader.Close()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sigchan:
			fmt.Println("Interrompendo o consumidor...")
			return
		default:
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Erro ao ler mensagem: %v", err)
				continue
			}

			var payload KafkaPayload
			err = json.Unmarshal(message.Value, &payload)
			if err != nil {
				log.Printf("Erro ao fazer parse da mensagem: %v", err)
				continue
			}

			var result *models.Product

			switch payload.Operation {
			case "create":
				result, err = products.Create(payload.Product, false)
			case "update":
				result, err = products.Update(payload.Product, false)
			default:
				log.Printf("Operação desconhecida: %s", payload.Operation)
				continue
			}

			if err != nil {
				log.Printf("Erro ao executar operação: %v", err)
				continue
			}

			log.Printf("Operação %s realizada com sucesso: %+v", payload.Operation, result)
		}
	}
}
