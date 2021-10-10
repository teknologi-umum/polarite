package handlers

import (
	"errors"
	"net/http"
	"polarite/business/models"
	"polarite/repository"
	"polarite/resources"

	"github.com/allegro/bigcache/v3"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func (d *Dependency) AddPaste(c *fiber.Ctx) error {
	body := c.Body()

	conn, err := d.DB.Connx(c.Context())
	if err != nil {
		return err
	}

	// Check duplicates
	hash, err := resources.Hash(body)
	if err != nil {
		return err
	}

	dup, i, err := d.PasteController.ReadHashFromDB(conn, hash)
	if err != nil {
		return err
	}

	if dup {
		c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
		return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL + i.ID))
	}

	conn, err = d.DB.Connx(c.Context())
	if err != nil {
		return err
	}

	paste := models.Item{
		Paste: string(body),
		Hash:  hash,
		IP:    c.IP(),
		User:  c.Locals("user").(string),
	}
	data, err := d.PasteController.InsertPasteToDB(conn, paste)
	if err != nil {
		return err
	}

	err = d.PasteController.InsertPasteToCache(data)
	if err != nil {
		return err
	}

	conn, err = d.DB.Connx(c.Context())
	if err != nil {
		return err
	}
	defer d.updateCachedID(conn, data.ID)

	c.Set(fiber.HeaderContentType, fiber.MIMETextPlainCharsetUTF8)
	return c.Status(http.StatusCreated).Send([]byte(repository.BASE_URL + data.ID))
}

func (d *Dependency) updateCachedID(conn *sqlx.Conn, id string) error {
	ids, err := d.PasteController.ReadIDFromMemory()
	if err != nil && !errors.Is(err, bigcache.ErrEntryNotFound) {
		return err
	}

	if errors.Is(err, bigcache.ErrEntryNotFound) {
		pastes, err := d.PasteController.ReadIDFromDB(conn)
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
