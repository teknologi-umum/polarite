package handlers

import (
	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type Dependency struct {
	DB     *sqlx.DB
	Memory *bigcache.BigCache
	Cache  *redis.Client
}
