package controllers

import (
	"context"
	"polarite/business/models"
	"strings"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/allegro/bigcache/v3"
	"github.com/georgysavva/scany/sqlscan"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
)

type PasteControllerImpl struct {
	Cache  *redis.Client
	Memory *bigcache.BigCache
}

func (c *PasteControllerImpl) ReadItemFromCache(id string) (models.Item, error) {

	r, err := c.Cache.Get(context.Background(), "paste:"+id).Result()
	if err != nil {
		return models.Item{}, err
	}

	result := models.Item{
		ID:    id,
		Paste: r,
	}

	return result, nil
}

func (c *PasteControllerImpl) ReadItemFromDB(db *sqlx.Conn, id string) (models.Item, error) {
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

func (c *PasteControllerImpl) ReadIDFromDB(db *sqlx.Conn) ([]models.Item, error) {
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

func (c *PasteControllerImpl) ReadIDFromMemory() ([]string, error) {
	s, err := c.Memory.Get("ids")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(s), ","), nil
}

func (c *PasteControllerImpl) InsertPasteToDB(db *sqlx.Conn, body []byte) (models.Item, error) {
	defer db.Close()

	id, err := nanoid.New()
	if err != nil {
		return models.Item{}, err
	}

	creationTime := time.Now().Format(time.RFC3339)
	r, err := db.QueryContext(context.Background(), "INSERT INTO paste (id, content, created) VALUES ($1, $2, $3)", id, string(body), creationTime)
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

func (c *PasteControllerImpl) InsertPasteToCache(paste models.Item) error {
	_, err := c.Cache.SetEX(context.Background(), "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *PasteControllerImpl) UpdateIDListFromDB(pastes []models.Item) ([]string, error) {
	var temp []string
	for i := 0; i < len(pastes); i++ {
		temp = append(temp, pastes[i].ID)
	}

	s := strings.Join(temp, ",")
	err := c.Memory.Set("ids", []byte(s))
	if err != nil {
		return []string{}, err
	}

	return temp, nil
}

func (c *PasteControllerImpl) UpdateIDListFromCache(pastes []string, new string) (int, error) {
	pastes = append(pastes, new)
	s := strings.Join(pastes, ",")
	err := c.Memory.Set("ids", []byte(s))
	if err != nil {
		return 0, err
	}

	return len(pastes), nil
}
