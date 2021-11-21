package paste

import (
	"context"

	"github.com/allegro/bigcache/v3"
	"github.com/getsentry/sentry-go"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	Cache  *redis.Client
	Memory *bigcache.BigCache
	DB     *sqlx.DB
	Logger *sentry.Client
}

type PasteController interface {
	ReadItemFromCache(ctx context.Context, id string) (Item, error)
	ReadItemFromDB(ctx context.Context, id string) (Item, error)
	ReadIDFromDB(ctx context.Context) ([]Item, error)
	ReadIDFromMemory() ([]string, error)
	ReadHashFromDB(ctx context.Context, h string) (bool, Item, error)
	InsertPasteToDB(ctx context.Context, paste Item) (Item, error)
	InsertPasteToCache(ctx context.Context, paste Item) error
	UpdateIDListFromDB(pastes []Item) ([]string, error)
	UpdateIDListFromCache(pastes []string, new string) (int, error)
}
