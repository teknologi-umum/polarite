// Get a data from database

package controllers

import (
	"context"
	"polarite/business/models"
	"strings"

	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

func ReadItemFromCache(cache *redis.Client, id string) (models.Item, error) {
	r, err := cache.Get(context.Background(), "paste:"+id).Result()
	if err != nil {
		return models.Item{}, err
	}

	result := models.Item{
		ID:    id,
		Paste: r,
	}

	return result, nil
}

func ReadItemFromDB(db *sqlx.Conn, id string) (models.Item, error) {
	defer db.Close()

	r, err := db.QueryContext(context.Background(), "SELECT content FROM paste WHERE id = $1", id)
	if err != nil {
		return models.Item{}, err
	}
	defer r.Close()

	var result models.Item
	err = sqlscan.ScanOne(&result, r)
	if err != nil {
		return models.Item{}, err
	}

	return result, nil
}

func ReadIDFromDB(db *sqlx.Conn) ([]models.Item, error) {
	defer db.Close()

	r, err := db.QueryContext(context.Background(), "SELECT id FROM paste")
	if err != nil {
		return []models.Item{}, err
	}
	defer r.Close()

	var result []models.Item
	err = sqlscan.ScanAll(&result, r)
	if err != nil {
		return []models.Item{}, err
	}

	return result, nil
}

func ReadIDFromMemory(memory *bigcache.BigCache) ([]string, error) {
	s, err := memory.Get("ids")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(s), ","), nil
}
