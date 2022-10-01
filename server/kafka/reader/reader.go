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
		Brokers:     strings.Split(os.Getenv("BROKER_LIST"), ","),
		Topic:       "operations",
		GroupID:     "consumer:operations.created",
		Dialer:      dialer,
		StartOffset: kafka.LastOffset,
	})

	for {
		m, err := Reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Msg(err.Error())
			break
		}

		log.Info().Msgf("Messaged received: %v", m)

		msg := new(messages.OperationCreatedMsg)
		unmarshalErr := json.Unmarshal(m.Value, &msg)

		if unmarshalErr != nil {
			log.Error().Msg(err.Error())
			break
		}

		log.Info().Msgf("Message parsed: %v", msg)

		// record the operation
		db.DB.Create(&db.Operation{
			ID:            msg.Data.ID,
			DocumentID:    msg.Data.DocumentID,
			IsProcessed:   msg.Data.IsProcessed,
			Action:        msg.Data.Action,
			StartPosition: msg.Data.StartPosition,
			EndPosition:   msg.Data.EndPosition,
			Text:          msg.Data.Text,
			CreatedAt:     msg.Data.CreatedAt,
		})

		log.Info().Msgf("Operation record with ID: %v", msg.Data.ID)

		// modify the document
	}
}
