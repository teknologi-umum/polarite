// Insert a data to the database

package paste

import (
	"context"
	"database/sql"
	"polarite/resources"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
)

func (c *Dependency) InsertPasteToDB(ctx context.Context, paste Item) (Item, error) {
	conn, err := c.DB.Connx(ctx)
	if err != nil {
		return Item{}, err
	}
	defer conn.Close()

	id, err := nanoid.New()
	if err != nil {
		return Item{}, err
	}

	p, err := resources.CompressContent(paste.Paste)
	if err != nil {
		return Item{}, err
	}

	tx, err := conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return Item{}, err
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		id, p, paste.Hash, paste.IP, paste.User)
	if err != nil {
		tx.Rollback()
		return Item{}, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return Item{}, err
	}

	return Item{
		ID:        id,
		Paste:     paste.Paste,
		CreatedAt: time.Now(),
	}, nil
}

func (c *Dependency) InsertPasteToCache(ctx context.Context, paste Item) error {
	_, err := c.Cache.SetEX(ctx, "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
