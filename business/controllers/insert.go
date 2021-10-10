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

	creationTime := time.Now().Format(time.RFC3339)
	r, err := db.QueryContext(
		context.Background(),
		"INSERT INTO paste (id, content, hash, created, ip, user) VALUES (?, ?, ?, ?, ?, ?)",
		id, paste.Paste, paste.Hash, creationTime, paste.IP, paste.User)
	if err != nil {
		return models.Item{}, err
	}
	defer r.Close()

	t, err := time.Parse(time.RFC3339, creationTime)
	if err != nil {
		return models.Item{}, err
	}

	return models.Item{
		ID:        id,
		Paste:     paste.Paste,
		CreatedAt: t,
	}, nil
}

func (c *PasteControllerImpl) InsertPasteToCache(paste models.Item) error {
	_, err := c.Cache.SetEX(context.Background(), "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
