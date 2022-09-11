package logger

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogStruct struct {
	URL       string `json:"url"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  int64  `json:"duration"`
	Status    int    `json:"status"`
	Method    string `json:"method"`
}

func (l LogStruct) MarshalZerologObject(e *zerolog.Event) {
	e.
		Str("method", l.Method).
		Str("url", l.URL).
		Int("status", l.Status).
		Str("start_time", l.StartTime).
		Str("end_time", l.EndTime).
		Int64("duration", l.Duration)
}

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		t := time.Now().UTC()
		logData := LogStruct{
			URL:       c.OriginalURL(),
			StartTime: t.String(),
			Method:    string(c.Method()),
		}

		c.Next()

		logData.Status = c.Response().StatusCode()
		logData.EndTime = time.Now().String()
		logData.Duration = time.Since(t).Milliseconds()

		log.Info().EmbedObject(logData).Send()

		return nil
	}
}
