package handlers

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
)

// Redirect any incoming HTTP request into HTTPS.
func GoSecure(c *fiber.Ctx) error {
	if !c.Secure() && os.Getenv("ENVIRONMENT") == "production" {
		c.Set("Location", fmt.Sprintf("https://%s%s", c.Hostname(), c.OriginalURL()))
		return c.SendStatus(fiber.StatusPermanentRedirect)
	}
	return c.Next()
}
