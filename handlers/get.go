package handlers

import (
	"errors"
	"net/http"
	"polarite/business/models"
	h "polarite/platform/highlight"
	"polarite/repository"
	"polarite/resources"
	"strings"

	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

// Get route to find a paste by given ID (on path parameters).
func (d *Dependency) Get(c *fiber.Ctx) error {
	// Parse the URL param first
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).Send([]byte(repository.ErrNoID.Error()))
	}

	var qs models.QueryString
	err := c.QueryParser(&qs)
	if err != nil {
		return err
	}

	// Validate if the ID exists or not
	ids, err := d.PasteController.ReadIDFromMemory()
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		conn, err := d.DB.Connx(c.Context())
		if err != nil {
			return err
		}

		pastes, err := d.PasteController.ReadIDFromDB(conn)
		if err != nil {
			return err
		}

		ids, err = d.PasteController.UpdateIDListFromDB(pastes)
		if err != nil {
			return err
		}
	}

	idExists := resources.ValidateID(ids, id)
	if !idExists {
		return c.Status(http.StatusNotFound).Send([]byte(repository.ID_NOT_FOUND))
	}

	// Try to fetch from Redis first
	i, err := d.PasteController.ReadItemFromCache(id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Item not found on Redis, now try to fetch from DB
	if errors.Is(err, redis.Nil) {
		conn, err := d.DB.Connx(c.Context())
		if err != nil {
			return err
		}

		i, err = d.PasteController.ReadItemFromDB(conn, id)
		if err != nil {
			return err
		}

		err = d.PasteController.InsertPasteToCache(i)
		if err != nil {
			return err
		}
	}

	// we need to replace escaped newline to literal newline
	content := strings.Replace(i.Paste, `\n`, "\n", -1)

	if qs.Language != "" {
		highlighted, err := h.Highlight(content, qs.Language, qs.Theme, qs.LineNr)
		if err != nil {
			// they should still be able to get the plain text even if the highlighter is b0rked
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(http.StatusOK).Send([]byte(content))
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.Status(http.StatusOK).Send([]byte(highlighted))
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusOK).Send([]byte(content))
}
