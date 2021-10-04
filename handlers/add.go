package handlers

import (
	"polarite/business/controllers"
	"polarite/repository"

	"github.com/gofiber/fiber/v2"
)

func (d *Dependency) AddPaste(c *fiber.Ctx) error {
	body := c.Body()

	exceeded := controllers.ValidateSize(body)
	if exceeded {
		return repository.ErrBodyTooBig
	}
	conn, err := d.DB.Acquire(c.Context())
	if err != nil {
		return err
	}
	
	data, err := controllers.InsertPasteToDB(conn, body)
	if err != nil {
		return err
	}

	c.Set("Content-Type", "text/plain")
	return c.Status(200).Send([]byte(repository.BASE_URL + data.ID))
}
