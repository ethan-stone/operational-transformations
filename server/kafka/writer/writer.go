package writer

import (
	"crypto/tls"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var Writer *kafka.Writer

func Connect() {
	mechanism, err := scram.Mechanism(scram.SHA256, os.Getenv("KAFKA_USERNAME"), os.Getenv("KAFKA_PASSWORD"))

	if err != nil {
		log.Fatal().Msg("Failed to initialize Kafka scram mechanism")
	}

	dialer := &kafka.Dialer{
		SASLMechanism: mechanism,
		TLS:           &tls.Config{},
	}

	Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      strings.Split(os.Getenv("BROKER_LIST"), ","),
		Topic:        "operations",
		Dialer:       dialer,
		BatchTimeout: 100 * time.Millisecond,
	})
}
