package controllers

import (
	"log"
	"net/http"
	"os"
	"polarite/platform/logtail"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
)

// Override Fiber's default Errorhandler so we the users can't see anything.
// Then we log the error to Logtail and Sentry.
// Send the error to stdout if "ENVIRONMENT" env variable is not set to "production".

func ErrorHandler(logger *sentry.Client) func(c *fiber.Ctx, e error) error {
	return func(c *fiber.Ctx, e error) error {
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

		scope := sentry.NewScope()
		scope.SetContext("app:meta", map[string]interface{}{
			"original url": c.OriginalURL(),
			"method":       c.Method(),
			"body":         string(c.Body()),
		})
		_ = logger.CaptureException(e, &sentry.EventHint{OriginalException: e}, scope)

		if os.Getenv("ENVIRONMENT") != "production" {
			log.Println(e)
		}

		return c.
			Status(http.StatusInternalServerError).
			Send([]byte("An error occured on the server"))
	}
}
