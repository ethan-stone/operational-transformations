package main

import (
	"os"
	"os/signal"

	"github.com/ethan-stone/optra/server/db"
	"github.com/ethan-stone/optra/server/kafka/reader"
	"github.com/ethan-stone/optra/server/kafka/writer"
	"github.com/ethan-stone/optra/server/middleware/logger"
	"github.com/ethan-stone/optra/server/router/document"
	"github.com/ethan-stone/optra/server/router/operation"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog/log"
)

func ping(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "pong",
	})
}

func main() {
	db.Connect()
	writer.Connect()
	go reader.Connect() // run the consumer asynchronously

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

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)

	go func() {
		<-signalChannel
		log.Info().Msg("Server shutting down...")
		reader.Reader.Close()
		app.Shutdown()
		log.Info().Msg("Server shutdown.")
	}()

	app.Listen(":8080")
}
