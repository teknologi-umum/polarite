// Get a data from database

package controllers

import (
	"context"
	"database/sql"
	"errors"
	"polarite/business/models"
	"polarite/resources"
	"strings"

	"github.com/georgysavva/scany/sqlscan"
)

func (c *PasteControllerImpl) ReadItemFromCache(ctx context.Context, id string) (models.Item, error) {
	r, err := c.Cache.Get(ctx, "paste:"+id).Result()
	if err != nil {
		return models.Item{}, err
	}

	result := models.Item{
		ID:    id,
		Paste: []byte(r),
	}

	return result, nil
}

func (c *PasteControllerImpl) ReadItemFromDB(ctx context.Context, id string) (models.Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return models.Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryContext(ctx, "SELECT content FROM paste WHERE id = ?", id)
	if err != nil {
		return models.Item{}, err
	}
	defer r.Close()

	var result models.Item
	err = sqlscan.ScanOne(&result, r)
	if err != nil {
		return models.Item{}, err
	}

	p, err := resources.DecompressContent(result.Paste)
	if err != nil {
		return models.Item{}, err
	}

	return models.Item{
		Paste: p,
	}, nil
}

func (c *PasteControllerImpl) ReadIDFromDB(ctx context.Context) ([]models.Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return []models.Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryContext(ctx, "SELECT id FROM paste")
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

func (c *PasteControllerImpl) ReadHashFromDB(ctx context.Context, h string) (bool, models.Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return false, models.Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryContext(ctx, "SELECT id FROM paste WHERE hash = ?", h)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.Item{}, nil
		}
		return false, models.Item{}, err
	}
	defer r.Close()

	var item models.Item
	err = sqlscan.ScanOne(&item, r)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, models.Item{}, nil
		}
		return false, models.Item{}, err
	}
	return true, item, nil
}
