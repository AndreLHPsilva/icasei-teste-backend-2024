package products

import (
	"context"
	"ms-go/app/helpers"
	"ms-go/app/models"
	"ms-go/db"
	"net/http"
	"time"

	"ms-go/app/producers"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(data models.Product, isAPI bool) (*models.Product, error) {

	if data.ID == 0 {
		var max models.Product

		opts := options.FindOne()

		opts.SetSort(bson.D{{Key: "created_at", Value: -1}})

		db.Connection().FindOne(context.TODO(), bson.D{}, opts).Decode(&max)

		data.ID = max.ID + 1
	}

	if err := data.Validate(); err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusUnprocessableEntity}
	}

	data.CreatedAt = time.Now()
	data.UpdatedAt = data.CreatedAt

	if data.Stock == 0 {
    data.Stock = 0
	}

	_, err := db.Connection().InsertOne(context.TODO(), data)

	if err != nil {
		return nil, &helpers.GenericError{Msg: err.Error(), Code: http.StatusInternalServerError}
	}

	defer db.Disconnect()

	if isAPI {
		payload := producers.KafkaPayload{
			Product: producers.ProductWithoutTimestamps{
				ID:          data.ID,
				Name:        data.Name,
				Brand:       data.Brand,
				Price:       data.Price,
				Description: data.Description,
				Stock:       data.Stock,
			},
		}
	
		err := producers.ProducerKafka(payload)
		if err != nil {
			return nil, err // Trate o erro de acordo com sua lógica de aplicação
		}
	}

	return &data, nil
}
