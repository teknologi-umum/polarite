package handlers

import (
	"errors"
	"net/http"
	"polarite/business/controllers"
	"polarite/repository"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
)

func (d *Dependency) AddPaste(c *fiber.Ctx) error {
	body := c.Body()

	exceeded := controllers.ValidateSize(body)
	if exceeded {
		return c.Status(http.StatusBadRequest).Send([]byte(repository.ErrBodyTooBig.Error()))
	}

	conn, err := d.DB.Acquire(c.Context())
	if err != nil {
		return err
	}

	data, err := controllers.InsertPasteToDB(conn, body)
	if err != nil {
		return err
	}

	err = controllers.InsertPasteToCache(d.Cache, data)
	if err != nil {
		return err
	}

	defer updateCachedID(conn, d.Memory, data.ID)

	c.Set("Content-Type", "text/plain")
	return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL + data.ID))
}

func updateCachedID(conn *pgxpool.Conn, mem *bigcache.BigCache, id string) error {
	ids, err := controllers.ReadIDFromMemory(mem)
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		pastes, err := controllers.ReadIDFromDB(conn)
		if err != nil {
			return err
		}

		_, err = controllers.UpdateIDListFromDB(mem, pastes)
		if err != nil {
			return err
		}
	}

	_, err = controllers.UpdateIDListFromCache(mem, ids, id)
	if err != nil {
		return err
	}

	return nil
}
