package main

import (
	"github.com/ethan-stone/optra/server/db"
	"github.com/ethan-stone/optra/server/middleware/logger"
	"github.com/ethan-stone/optra/server/router/document"
	"github.com/ethan-stone/optra/server/router/operation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "pong",
	})
}

func main() {
	db.Connect()

	app := fiber.New()

	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowHeaders:     "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With",
	}))

	app.Get("/ping", ping)

	app.Post("/documents", document.Create)
	app.Post("/operations", operation.Create)

	app.Listen(":8080")
}
