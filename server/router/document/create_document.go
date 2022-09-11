package document

import (
	"github.com/ethan-stone/optra/server/db"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

func Create(c *fiber.Ctx) error {
	document := new(db.Document)
	result := db.DB.Create(&document)

	if result.Error != nil {
		log.Error().Msg(result.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	log.Info().Msgf("Document With ID: %v created", document.ID)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": document,
	})
}
