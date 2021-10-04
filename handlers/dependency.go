package handlers

import (
	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Dependency struct {
	DB     *pgxpool.Pool
	Memory *bigcache.BigCache
	Cache  *redis.Client
}
