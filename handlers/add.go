package handlers

import (
	"context"
	"errors"
	"net/http"
	"polarite/packages/paste"
	"polarite/repository"
	"polarite/resources"

	"github.com/allegro/bigcache/v3"
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

	dup, i, err := d.PasteController.ReadHashFromDB(c.Context(), hash)
	if err != nil {
		return err
	}

	if dup {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL() + i.ID))
	}

	item := paste.Item{
		Paste: body,
		Hash:  hash,
		IP:    c.IP(),
		User:  c.Locals("user").(string),
	}
	data, err := d.PasteController.InsertPasteToDB(c.Context(), item)
	if err != nil {
		return err
	}

	err = d.PasteController.InsertPasteToCache(c.Context(), data)
	if err != nil {
		return err
	}

	defer d.updateCachedID(data.ID)

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL() + data.ID))
}

// A private function to run without wanting to wait.
// This will return behind the scene once, preferably without
// the use of goroutine.
func (d *Dependency) updateCachedID(id string) error {
	ids, err := d.PasteController.ReadIDFromMemory()
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		pastes, err := d.PasteController.ReadIDFromDB(context.Background())
		if err != nil {
			return err
		}

		_, err = d.PasteController.UpdateIDListFromDB(pastes)
		if err != nil {
			return err
		}
	}

	_, err = d.PasteController.UpdateIDListFromCache(ids, id)
	if err != nil {
		return err
	}

	return nil
}
