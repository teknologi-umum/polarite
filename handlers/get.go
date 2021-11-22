package handlers

import (
	"errors"
	"net/http"
	"polarite/packages/paste"
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

	var qs paste.QueryString
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
		pastes, err := d.PasteController.ReadIDFromDB(c.Context())
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
	i, err := d.PasteController.ReadItemFromCache(c.Context(), id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	// Item not found on Redis, now try to fetch from DB
	if errors.Is(err, redis.Nil) {
		i, err = d.PasteController.ReadItemFromDB(c.Context(), id)
		if err != nil {
			return err
		}

		err = d.PasteController.InsertPasteToCache(c.Context(), i)
		if err != nil {
			return err
		}
	}

	// we need to replace escaped newline to literal newline
	content := strings.Replace(string(i.Paste), `\n`, "\n", -1)

	if qs.Language != "" {
		highlighted, err := h.Highlight(content, qs.Language, qs.Theme, qs.LineNr)
		if err != nil {
			// they should still be able to get the plain text even if the highlighter is b0rked
			c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
			return c.Status(http.StatusOK).Send([]byte(content))
		}

		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.Render("template", fiber.Map{
			"Content":     "",
			"Highlighted": highlighted,
			"Truncated":   resources.TruncateString(content, 50),
		})
	}

	if strings.Contains(c.Get(fiber.HeaderAccept), fiber.MIMETextHTML) {
		return c.Render("template", fiber.Map{
			"Content":     content,
			"Highlighted": "",
			"Truncated":   resources.TruncateString(content, 50),
		})
	}

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusOK).Send([]byte(content))
}
