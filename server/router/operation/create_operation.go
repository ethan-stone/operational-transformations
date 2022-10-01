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
	DocumentID    string  `json:"document_id" binding:"required"`
	Action        string  `json:"action" binding:"required"`
	StartPosition int32   `json:"start_position" binding:"required"`
	EndPosition   *int32  `json:"end_position,omitempty"`
	Text          *string `json:"text,omitempty"`
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

	if body.Action == "INSERT" && body.Text == nil {
		log.Error().Msgf("INSERT operation missing text, BadRequest")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "INSERT operation is missing text",
		})
	}

	if body.Action == "DELETE" && body.EndPosition == nil {
		log.Error().Msgf("DELETE operation missing end position, BadRequest")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "DELETE operation is missing end position",
		})
	}

	operationCreatedMsgData := messages.OperationCreatedMsgData{
		ID:            uuid.New(),
		DocumentID:    documentId,
		IsProcessed:   false,
		Action:        body.Action,
		StartPosition: body.StartPosition,
		CreatedAt:     time.Now(),
	}

	if body.Action == "DELETE" {
		operationCreatedMsgData.EndPosition = *body.EndPosition
		operationCreatedMsgData.Text = ""
	}

	if body.Action == "INSERT" {
		operationCreatedMsgData.EndPosition = body.StartPosition
		operationCreatedMsgData.Text = *body.Text
	}

	msgJSON := messages.NewOperationCreatedMsg(operationCreatedMsgData)

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
