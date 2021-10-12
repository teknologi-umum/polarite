package controllers_test

import (
	"context"
	"polarite/business/controllers"
	"polarite/business/models"
	"testing"
)

func TestReadItemFromCache(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	err := rds.Set(context.Background(), "paste:testid", "Hello world!", 0).Err()
	if err != nil {
		t.Error(err)
	}

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.ReadItemFromCache("testid")
	if err != nil {
		t.Error(err)
	}

	if i.Paste != "Hello world!" {
		t.Error("i.Paste should be \"Hello world!\", got:", i.Paste)
	}
}

func TestReadItemFromDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	c, err := db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	paste := models.Item{
		ID:    "wNnwj138ne9ZaWmNADwIg",
		Paste: "Hello world!",
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	r, err := c.QueryContext(
		context.Background(),
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		paste.ID,
		paste.Paste,
		paste.Hash,
		paste.IP,
		paste.User,
	)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	c, err = db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.ReadItemFromDB(c, paste.ID)
	if err != nil {
		t.Error(err)
	}

	if i.Paste != paste.Paste {
		t.Error("i.Paste should be "+paste.Paste+", got:", i.Paste)
	}
}

func TestReadIDFromDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	c, err := db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	paste := models.Item{
		ID:    "wNnwj138ne9ZaWmNADwIg",
		Paste: "Hello world!",
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	r, err := c.QueryContext(
		context.Background(),
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		paste.ID,
		paste.Paste,
		paste.Hash,
		paste.IP,
		paste.User,
	)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	c, err = db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.ReadIDFromDB(c)
	if err != nil {
		t.Error(err)
	}

	if len(i) != 1 {
		t.Error("length of i should be 1, got:", len(i))
	}

	if i[0].ID != paste.ID {
		t.Error("i[0].ID should be equal to "+paste.ID+", got:", i[0].ID)
	}
}

func TestReadIDFromMemory(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	err := mem.Set("ids", []byte("a,b,c,d,e"))
	if err != nil {
		t.Error(err)
	}

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.ReadIDFromMemory()
	if err != nil {
		t.Error(err)
	}

	if len(i) != 5 {
		t.Error("length of i should be 5, got:", len(i))
	}
}

func TestReadHashFromDB_Dup(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	c, err := db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	paste := models.Item{
		ID:    "wNnwj138ne9ZaWmNADwIg",
		Paste: "Hello world!",
		Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
		IP:    "127.0.0.1",
		User:  "example@test.com",
	}

	r, err := c.QueryContext(
		context.Background(),
		"INSERT INTO paste (id, content, hash, ip, user) VALUES (?, ?, ?, ?, ?)",
		paste.ID,
		paste.Paste,
		paste.Hash,
		paste.IP,
		paste.User,
	)
	if err != nil {
		t.Error(err)
	}
	defer r.Close()

	c, err = db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	b, i, err := p.ReadHashFromDB(c, paste.Hash)
	if err != nil {
		t.Error(err)
	}

	if !b {
		t.Error("b should be true, got:", b)
	}

	if i.ID != paste.ID {
		t.Error("i.ID should be equal to "+paste.ID+", got:", i.ID)
	}
}

func TestReadHashFromDB_NoDup(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	c, err := db.Connx(context.Background())
	if err != nil {
		t.Error(err)
	}
	defer c.Close()

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	b, _, err := p.ReadHashFromDB(c, "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4")
	if err != nil {
		t.Error(err)
	}

	if b {
		t.Error("b should be false, got:", b)
	}
}
