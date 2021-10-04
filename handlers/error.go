package handlers

import (
	"net/http"
	"polarite/platform/logtail"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, e error) error {
	err := logtail.Log(logtail.Error{
		Meta: logtail.Meta{
			Level: "error",
			Date: time.Now().String(),
		},
		Error: e.Error(),
	})
	if err != nil {
		return err
	}

	return c.
		Status(http.StatusInternalServerError).
		Send([]byte("An error occured on the server"))
}