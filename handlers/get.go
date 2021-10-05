package handlers

import (
	"errors"
	"net/http"
	"polarite/business/controllers"
	"polarite/business/models"
	"polarite/repository"

	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
)

func (d *Dependency) Get(c *fiber.Ctx) error {
	// Parse the URL param first
	id := c.Params("id")
	if id == "" {
		return c.Status(http.StatusBadRequest).Send([]byte(repository.ErrNoID.Error()))
	}

	// TODO: Process the query string
	var qs models.QueryString
	err := c.QueryParser(&qs)
	if err != nil {
		return err
	}

	conn, err := d.DB.Acquire(c.Context())
	if err != nil {
		return err
	}

	// Validate if the ID exists or not
	ids, err := controllers.ReadIDFromMemory(d.Memory)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		pastes, err := controllers.ReadIDFromDB(conn)
		if err != nil {
			return err
		}

		ids, err = controllers.UpdateIDListFromDB(d.Memory, pastes)
		if err != nil {
			return err
		}
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

	// Item not found on Redis, now try to fetch from DB
	if errors.Is(err, redis.Nil) {
		i, err = controllers.ReadItemFromDB(conn, id)
		if err != nil {
			return err
		}

		err = controllers.InsertPasteToCache(d.Cache, i)
		if err != nil {
			return err
		}
	}

	c.Set("Content-Type", "text/plain")
	return c.Status(http.StatusOK).Send([]byte(i.Paste))
}
