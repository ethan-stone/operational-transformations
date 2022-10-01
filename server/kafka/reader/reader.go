package reader

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"os"
	"strings"

	"github.com/ethan-stone/optra/server/db"
	"github.com/ethan-stone/optra/server/kafka/messages"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var Reader *kafka.Reader

func Connect() {
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"))

	if err != nil {
		log.Fatal().Msg("Failed to initialize Kafka scram mechanism")
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: strings.Split(os.Getenv("BROKER_LIST"), ","),
		Topic:   "operations",
		Dialer:  dialer,
	})

	for {
		m, err := Reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Msg(err.Error())
			break
		}
		msg := new(messages.OperationCreatedMsg)
		unmarshalErr := json.Unmarshal(m.Value, &msg)

		if unmarshalErr != nil {
			log.Error().Msg(err.Error())
			break
		}

		// record the operation
		db.DB.Create(&db.Operation{
			ID:          msg.Data.ID,
			DocumentID:  msg.Data.DocumentID,
			IsProcessed: msg.Data.IsProcessed,
			CreatedAt:   msg.Data.CreatedAt,
		})

		// modify the document
	}
}
