package paste_test

import (
	"context"
	"polarite/packages/paste"
	"testing"
)

func TestInsertPasteToDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := paste.Dependency{
		Cache:  rds,
		Memory: mem,
		DB:     db,
	}

	paste := paste.Item{
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
