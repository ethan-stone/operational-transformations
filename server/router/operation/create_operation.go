package operation

import (
	"github.com/ethan-stone/optra/server/db"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
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

	operation := new(db.Operation)
	operation.DocumentID = documentId
	operation.IsProcessed = false
	result := db.DB.Create(&operation)

	if result.Error != nil {
		log.Error().Msg(result.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	log.Info().Msgf("Operation With ID: %v created", operation.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": operation,
	})
}
