// Insert a data to the database

package controllers

import (
	"context"
	"polarite/business/models"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertPasteToDB(db *pgxpool.Conn, body []byte) (models.Item, error) {
	defer db.Release()

	id, err := nanoid.New()
	if err != nil {
		return models.Item{}, err
	}

	creationTime := time.Now().Format(time.RFC3339)
	r, err := db.Query(context.Background(), "INSERT INTO paste (id, content, created) VALUES ($1, $2, $3)", id, string(body), creationTime)
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
		Paste:     string(body),
		CreatedAt: t,
	}, nil
}

func InsertPasteToCache(cache *redis.Client, paste models.Item) error {
	_, err := cache.SetEX(context.Background(), "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
