package controllers

import (
	"net/http"

	"polarite/repository"
	"polarite/resources"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/gofiber/fiber/v2"
)

// Post Route to add a paste from somewhere.
// If the submitted paste is a duplicate, it will quickly return
// a generated ID based on the SHA224 hash.
func (d *Dependency) AddPaste(c *fiber.Ctx) error {
	body := c.Body()

	// Check duplicates
	hash, err := resources.Hash(body)
	if err != nil {
		return err
	}

	exists, existingId, err := d.Paste.ReadHash(c.Context(), hash)
	if err != nil {
		return err
	}

	if exists {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(http.StatusCreated).Send([]byte(repository.BaseUrl() + existingId))
	}

	id, err := nanoid.GenerateString(nanoid.DefaultAlphabet, 6)
	if err != nil {
		return err
	}

	item := repository.Item{
		ID:    id,
		Paste: body,
		Hash:  hash,
	}

	data, err := d.Paste.InsertPaste(c.Context(), item)
	if err != nil {
		return err
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusCreated).Send([]byte(repository.BaseUrl() + data.ID))
}
