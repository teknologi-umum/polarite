// Insert a data to the database

package controllers

import (
	"context"
	"polarite/business/models"
	"polarite/resources"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
)

func (c *PasteControllerImpl) InsertPasteToDB(ctx context.Context, paste models.Item) (models.Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return models.Item{}, err
	}
	defer conn.Close()

	id, err := nanoid.New()
	if err != nil {
		return models.Item{}, err
	}

	p, err := resources.CompressContent(paste.Paste)
	if err != nil {
		return models.Item{}, err
	}

	r, err := conn.QueryContext(
		ctx,
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		id, p, paste.Hash, paste.IP, paste.User)
	if err != nil {
		return models.Item{}, err
	}
	defer r.Close()

	if err != nil {
		return models.Item{}, err
	}

	return models.Item{
		ID:        id,
		Paste:     paste.Paste,
		CreatedAt: time.Now(),
	}, nil
}

func (c *PasteControllerImpl) InsertPasteToCache(ctx context.Context, paste models.Item) error {
	_, err := c.Cache.SetEX(ctx, "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
