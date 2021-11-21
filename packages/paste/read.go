// Get a data from database

package paste

import (
	"context"
	"database/sql"
	"errors"
	"polarite/resources"
	"strings"
)

func (c *Dependency) ReadItemFromCache(ctx context.Context, id string) (Item, error) {
	r, err := c.Cache.Get(ctx, "paste:"+id).Result()
	if err != nil {
		return Item{}, err
	}

	result := Item{
		ID:    id,
		Paste: []byte(r),
	}

	return result, nil
}

func (c *Dependency) ReadItemFromDB(ctx context.Context, id string) (Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryxContext(ctx, "SELECT content FROM paste WHERE id = ?", id)
	if err != nil {
		return Item{}, err
	}
	defer r.Close()

	var result Item
	for r.Next() {
		err = r.StructScan(&result)
		if err != nil {
			return Item{}, err
		}
	}

	p, err := resources.DecompressContent(result.Paste)
	if err != nil {
		return Item{}, err
	}

	return Item{
		Paste: p,
	}, nil
}

func (c *Dependency) ReadIDFromDB(ctx context.Context) ([]Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return []Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryxContext(ctx, "SELECT id FROM paste")
	if err != nil {
		return []Item{}, err
	}
	defer r.Close()

	var result []Item
	for r.Next() {
		var item Item
		err = r.StructScan(&item)
		if err != nil {
			return []Item{}, err
		}
		result = append(result, item)
	}

	return result, nil
}

func (c *Dependency) ReadIDFromMemory() ([]string, error) {
	s, err := c.Memory.Get("ids")
	if err != nil {
		return []string{}, err
	}

	return strings.Split(string(s), ","), nil
}

func (c *Dependency) ReadHashFromDB(ctx context.Context, h string) (bool, Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return false, Item{}, err
	}
	defer conn.Close()

	r, err := conn.QueryxContext(ctx, "SELECT id FROM paste WHERE hash = ?", h)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, Item{}, nil
		}
		return false, Item{}, err
	}
	defer r.Close()

	var item Item
	for r.Next() {
		err = r.StructScan(&item)
		if err != nil {
			return false, Item{}, err
		}
	}

	if item.ID == "" {
		return false, Item{}, nil
	}

	return true, item, nil
}
