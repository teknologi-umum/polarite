package handlers

import (
	"log"
	"net/http"
	"os"
	"polarite/platform/logtail"
	"time"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, e error) error {
	if e.Error() == "Method Not Allowed" {
		return c.Status(http.StatusMethodNotAllowed).Send([]byte(e.Error()))
	}

	err := logtail.Log(logtail.Error{
		Meta: logtail.Meta{
			Level: "error",
			Date:  time.Now().String(),
		},
		Error: e.Error(),
	})
	if err != nil {
		return err
	}

	if os.Getenv("ENVIRONMENT") != "production" {
		log.Println(e)
	}

	return c.
		Status(http.StatusInternalServerError).
		Send([]byte("An error occured on the server"))
}
