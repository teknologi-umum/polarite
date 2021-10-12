// Insert a data to the database

package controllers

import (
	"context"
	"polarite/business/models"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/jmoiron/sqlx"
)

func (c *PasteControllerImpl) InsertPasteToDB(db *sqlx.Conn, paste models.Item) (models.Item, error) {
	defer db.Close()

	id, err := nanoid.New()
	if err != nil {
		return models.Item{}, err
	}

	r, err := db.QueryContext(
		context.Background(),
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		id, paste.Paste, paste.Hash, paste.IP, paste.User)
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

func (c *PasteControllerImpl) InsertPasteToCache(paste models.Item) error {
	_, err := c.Cache.SetEX(context.Background(), "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
