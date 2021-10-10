package handlers

import (
	"errors"
	"net/http"
	"polarite/business/controllers"
	"polarite/repository"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func (d *Dependency) AddPaste(c *fiber.Ctx) error {
	body := c.Body()

	exceeded := controllers.ValidateSize(body)
	if exceeded {
		return c.Status(http.StatusBadRequest).Send([]byte(repository.ErrBodyTooBig.Error()))
	}

	conn, err := d.DB.Connx(c.Context())
	if err != nil {
		return err
	}

	data, err := controllers.InsertPasteToDB(conn, body)
	if err != nil {
		return err
	}

	err = d.PasteController.InsertPasteToCache(data)
	if err != nil {
		return err
	}

	defer d.updateCachedID(conn, data.ID)

	c.Set("Content-Type", "text/plain")
	return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL + data.ID))
}

func (d *Dependency) updateCachedID(conn *sqlx.Conn, id string) error {
	ids, err := d.PasteController.ReadIDFromMemory()
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		pastes, err := controllers.ReadIDFromDB(conn)
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
