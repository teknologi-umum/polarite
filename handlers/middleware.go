package handlers

import (
	"net/http"
	"polarite/repository"
	"polarite/resources"

	"github.com/gofiber/fiber/v2"
)

func ValidateInput(c *fiber.Ctx) error {
	body := c.Body()

	exceeded := resources.ValidateSize(body)
	if exceeded {
		return c.Status(http.StatusBadRequest).Send([]byte(repository.ErrBodyTooBig.Error()))
	}

	auth := c.Get(fiber.HeaderAuthorization)

	if auth == "" || len(auth) < 15 {
		return c.Status(http.StatusUnauthorized).Send([]byte(repository.ErrNoAuthHeader.Error()))
	}

	if len(auth) > 250 {
		c.Locals("user", auth[0:250])
	} else {
		c.Locals("user", auth)
	}

	return c.Next()
}
