package controllers

import (
	"polarite/business/models"

	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type PasteControllerImpl struct {
	Cache  *redis.Client
	Memory *bigcache.BigCache
}

type PasteController interface {
	ReadItemFromCache(id string) (models.Item, error)
	ReadItemFromDB(db *sqlx.Conn, id string) (models.Item, error)
	ReadIDFromDB(db *sqlx.Conn) ([]models.Item, error)
	ReadIDFromMemory() ([]string, error)
	ReadHashFromDB(db *sqlx.Conn, h string) (bool, models.Item, error)
	InsertPasteToDB(db *sqlx.Conn, paste models.Item) (models.Item, error)
	InsertPasteToCache(paste models.Item) error
	UpdateIDListFromDB(pastes []models.Item) ([]string, error)
	UpdateIDListFromCache(pastes []string, new string) (int, error)
}
