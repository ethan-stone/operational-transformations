package operation

import (
	"context"
	"encoding/json"
	"time"

	"github.com/ethan-stone/optra/server/kafka/messages"
	"github.com/ethan-stone/optra/server/kafka/writer"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
)

type CreateOperationInput struct {
	DocumentID string `json:"document_id" binding:"required"`
}

func Create(c *fiber.Ctx) error {
	body := new(CreateOperationInput)

	if err := c.BodyParser(body); err != nil {
		log.Error().Msg(err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	documentId, documentIdErr := uuid.Parse(body.DocumentID)

	if documentIdErr != nil {
		log.Error().Msg(documentIdErr.Error())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid Document ID format",
		})
	}

	msgJSON := messages.NewOperationCreatedMsg(messages.OperationCreatedMsgData{
		ID:          uuid.New(),
		DocumentID:  documentId,
		IsProcessed: false,
		CreatedAt:   time.Now(),
	})

	msg, msgMarshalErr := json.Marshal(msgJSON)

	if msgMarshalErr != nil {
		log.Error().Msg(msgMarshalErr.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	writer.Writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(documentId.String()),
		Value: msg,
	})

	log.Info().Msgf("OperationCreatedMsg sent with ID: %v created", msgJSON.ID)
	return c.SendStatus(fiber.StatusNoContent)
}
