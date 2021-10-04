package handlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (d *Dependency) HomePage(c *fiber.Ctx) error {
	// serves the home page
	return c.Status(http.StatusOK).Send([]byte("Coming soon"))
}
