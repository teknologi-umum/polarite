package handlers

import (
	"errors"
	"net/http"
	"polarite/business/controllers"
	"polarite/business/models"
	h "polarite/platform/highlight"
	"polarite/repository"
	"strings"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (d *Dependency) Get(c *fiber.Ctx) error {
	// Parse the URL param first
	id := c.Params("id")
	if id == "" {
		return repository.ErrNoID
	}

	// TODO: Process the query string
	var qs models.QueryString
	err := c.QueryParser(&qs)
	if err != nil {
		return err
	}

	// Validate if the ID exists or not
	ids, err := controllers.ReadIDFromMemory(d.Memory)
	if err != nil {
		return err
	}

	idExists := controllers.ValidateID(ids, id)
	if !idExists {
		return c.Status(http.StatusNotFound).Send([]byte(repository.ID_NOT_FOUND))
	}

	// Try to fetch from Redis first
	i, err := controllers.ReadItemFromCache(d.Cache, id)
	if err != nil && !errors.Is(err, redis.Nil) {
		return err
	}

	if errors.Is(err, redis.Nil) {
		conn, err := d.DB.Acquire(c.Context())
		if err != nil {
			return err
		}
		i, err = controllers.ReadItemFromDB(conn, id)
		if err != nil {
			return err
		}
	}

	// we need to replace escaped newline to literal newline
	content := strings.Replace(i.Paste, `\n`, "\n", -1)

	c.Set("Content-Type", "text/plain")
	if qs.Language != "" {
		highlighted, err := h.Highlight(content, qs.Language, qs.Theme, qs.LineNr)
		if err != nil {
			// they should still be able to get the plain text even if the highlighter is b0rked
			c.Status(http.StatusOK).Send([]byte(content))

			return err
		}

		c.Set("Content-Type", "text/html")
		return c.Status(http.StatusOK).Send([]byte(highlighted))
	}
	return c.Status(http.StatusOK).Send([]byte(content))
}
