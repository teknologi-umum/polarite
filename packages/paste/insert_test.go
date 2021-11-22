package paste_test

import (
	"context"
	"database/sql"
	"errors"
	"polarite/packages/paste"
	"testing"
	"time"
)

func TestInsertPasteToDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := paste.Dependency{
		Cache:  rds,
		Memory: mem,
		DB:     db,
	}

	paste := paste.Item{
		ID:    "asdfgh",
		Paste: []byte("Hello world!"),
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	i, err := p.InsertPasteToDB(context.Background(), paste)
	if err != nil {
		t.Error(err)
	}

	if i.ID == "" {
		t.Error("something went wrong, got:", i)
	}
}

func TestInsertPasteToDB_Collision(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(time.Second*30))
	defer cancel()
	c, err := db.Connx(ctx)
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	item := paste.Item{
		ID:    "asdfgh",
		Paste: []byte("Hello world!"),
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c63jk",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	tx, err := c.BeginTxx(ctx, &sql.TxOptions{})
	if err != nil {
		t.Error(err)
	}

	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		item.ID, item.Paste, item.Hash, item.IP, item.User)
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		t.Error(err)
	}

	p := paste.Dependency{
		Cache:  rds,
		Memory: mem,
		DB:     db,
	}

	i, err := p.InsertPasteToDB(context.Background(), item)
	if err == nil {
		t.Error("expected error, got:", i)
	}

	if err != nil && !errors.Is(err, paste.ErrIDDuplicate) {
		t.Error(err)
	}
}

func TestInsertPasteToCache(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := paste.Dependency{
		Cache:  rds,
		Memory: mem,
		DB:     db,
	}

	paste := paste.Item{
		ID:    "wNnwj138ne9ZaWmNADwIg",
		Paste: []byte("Hello world!"),
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	err := p.InsertPasteToCache(context.Background(), paste)
	if err != nil {
		t.Error(err)
	}
}
