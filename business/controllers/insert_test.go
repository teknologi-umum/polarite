package controllers_test

import (
	"context"
	"polarite/business/controllers"
	"polarite/business/models"
	"testing"
)

func TestInsertPasteToDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	c, err := db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	paste := models.Item{
		Paste: []byte("Hello world!"),
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	i, err := p.InsertPasteToDB(c, paste)
	if err != nil {
		t.Error(err)
	}

	if i.ID == "" {
		t.Error("something went wrong, got:", i)
	}
}

func TestInsertPasteToCache(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	paste := models.Item{
		ID:    "wNnwj138ne9ZaWmNADwIg",
		Paste: []byte("Hello world!"),
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	err := p.InsertPasteToCache(paste)
	if err != nil {
		t.Error(err)
	}
}
