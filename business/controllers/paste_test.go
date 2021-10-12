package controllers_test

import (
	"context"
	"os"
	"polarite/resources"
	"testing"
	"time"

	"github.com/allegro/bigcache/v3"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var db *sqlx.DB
var rds *redis.Client
var mem *bigcache.BigCache

func TestMain(m *testing.M) {
	Setup()
	defer Teardown()

	os.Exit(m.Run())
}

func Setup() {
	dbURL, err := resources.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	db, err = sqlx.Connect("mysql", dbURL)
	if err != nil {
		panic(err)
	}

	// Setup Redis
	rdsConfig, err := redis.ParseURL(os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	rds = redis.NewClient(rdsConfig)

	// Setup In-Memory
	mem, err = bigcache.NewBigCache(bigcache.DefaultConfig(time.Hour * 12))
	if err != nil {
		panic(err)
	}

	c, err := db.Conn(context.Background())
	if err != nil {
		panic(err)
	}
	defer c.Close()

	r, err := c.QueryContext(context.Background(), `CREATE TABLE IF NOT EXISTS paste (
		id VARCHAR(36) PRIMARY KEY NOT NULL,
		content MEDIUMTEXT NOT NULL,
		created TIMESTAMP NOT NULL DEFAULT NOW(),
		hash VARCHAR(255) UNIQUE NOT NULL,
		ip VARCHAR(20) NOT NULL,
		user VARCHAR(255) NOT NULL
	);`)
	if err != nil {
		panic(err)
	}
	defer r.Close()

}

func TruncateTable(db *sqlx.DB, rds *redis.Client, mem *bigcache.BigCache) error {
	c, err := db.Connx(context.Background())
	if err != nil {
		return err
	}
	defer c.Close()

	r, err := c.QueryContext(context.Background(), "TRUNCATE TABLE paste;")
	if err != nil {
		return err
	}
	defer r.Close()

	err = rds.FlushAll(context.Background()).Err()
	if err != nil {
		return err
	}

	err = mem.Reset()
	if err != nil {
		return err
	}

	return nil
}

func Teardown() {
	db.Close()
	mem.Close()
	rds.Close()
}
