// Insert a data to the database

package paste

import (
	"context"
	"database/sql"
	"errors"
	"polarite/resources"
	"time"
)

func (d *Dependency) InsertPasteToDB(ctx context.Context, paste Item) (Item, error) {
	conn, err := d.DB.Connx(ctx)
	if err != nil {
		return Item{}, err
	}
	defer conn.Close()

	p, err := resources.CompressContent(paste.Paste)
	if err != nil {
		return Item{}, err
	}

	tx, err := conn.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		return Item{}, err
	}

	// Make sure the id does not exists in the first place
	r, err := tx.QueryxContext(
		ctx,
		"SELECT id FROM paste WHERE id = ?",
		paste.ID,
	)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		tx.Rollback()
		return Item{}, err
	}

	if r.Next() {
		tx.Rollback()
		return Item{}, ErrIDDuplicate
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		paste.ID, p, paste.Hash, paste.IP, paste.User)
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
		ID:        paste.ID,
		Paste:     paste.Paste,
		CreatedAt: time.Now(),
	}, nil
}

func (d *Dependency) InsertPasteToCache(ctx context.Context, paste Item) error {
	_, err := d.Cache.SetEX(ctx, "paste:"+paste.ID, paste.Paste, time.Hour*24*2).Result()
	if err != nil {
		return err
	}

	return nil
}
