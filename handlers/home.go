package handlers

import (
	"github.com/gofiber/fiber/v2"
)

// Home page is the index of GET /.
// This is supposed to be interactive page that users can store their paste
// and get the corresponding link.
func (d *Dependency) HomePage(c *fiber.Ctx) error {
	// serves the home page
	return c.Render("home", nil)
}
