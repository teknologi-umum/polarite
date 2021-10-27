package controllers

import (
	"context"
	"polarite/business/models"

	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type PasteControllerImpl struct {
	Cache  *redis.Client
	Memory *bigcache.BigCache
	DB     *sqlx.DB
}

type PasteController interface {
	ReadItemFromCache(ctx context.Context, id string) (models.Item, error)
	ReadItemFromDB(ctx context.Context, id string) (models.Item, error)
	ReadIDFromDB(ctx context.Context) ([]models.Item, error)
	ReadIDFromMemory() ([]string, error)
	ReadHashFromDB(ctx context.Context, h string) (bool, models.Item, error)
	InsertPasteToDB(ctx context.Context, paste models.Item) (models.Item, error)
	InsertPasteToCache(ctx context.Context, paste models.Item) error
	UpdateIDListFromDB(pastes []models.Item) ([]string, error)
	UpdateIDListFromCache(pastes []string, new string) (int, error)
}
