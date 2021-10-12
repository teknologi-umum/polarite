package controllers_test

import (
	"polarite/business/controllers"
	"polarite/business/models"
	"strings"
	"testing"
)

func TestUpdateIDListFromDB(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	paste := []models.Item{
		{
			ID:    "wNnwj138ne9ZaWmNADwIg",
			Paste: "Hello world!",
			Hash:  "7e81ebe9e604a0c97fef0e4cfe71f9ba0ecba13332bde953ad1c66e4",
			IP:    "127.0.0.1",
			User:  "example@test.com",
		},
		{
			ID:    "b_ZbHoI3gQTv4mR0CDLNQ",
			Paste: "Java sucks",
			Hash:  "da1d7ce7e6bdc6f5c88b448afbb0d14afaa338d0f0b6f85c02b451e2",
			IP:    "127.0.0.1",
			User:  "example@test.com",
		},
	}

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.UpdateIDListFromDB(paste)
	if err != nil {
		t.Error(err)
	}

	if len(i) != 2 {
		t.Error("length of i should be 2, got:", 2)
	}

	if strings.Join(i, ",") != "wNnwj138ne9ZaWmNADwIg,b_ZbHoI3gQTv4mR0CDLNQ" {
		t.Error("i should equal to []string{wNnwj138ne9ZaWmNADwIg,b_ZbHoI3gQTv4mR0CDLNQ}, got:", strings.Join(i, ","))
	}

	m, err := mem.Get("ids")
	if err != nil {
		t.Error(err)
	}

	if string(m) != "wNnwj138ne9ZaWmNADwIg,b_ZbHoI3gQTv4mR0CDLNQ" {
		t.Error("m is not equal to what expected (it's a bit long), got:", string(m))
	}
}

func TestUpdateIDListFromCache(t *testing.T) {
	defer TruncateTable(db, rds, mem)

	p := controllers.PasteControllerImpl{
		Cache:  rds,
		Memory: mem,
	}

	i, err := p.UpdateIDListFromCache([]string{"wNnwj138ne9ZaWmNADwIg", "b_ZbHoI3gQTv4mR0CDLNQ"}, "DQNYSjRH7AYthVmJ7P9-T")
	if err != nil {
		t.Error(err)
	}

	if i != 3 {
		t.Error("i should be equal to 3, got:", i)
	}

	m, err := mem.Get("ids")
	if err != nil {
		t.Error(err)
	}

	if string(m) != "wNnwj138ne9ZaWmNADwIg,b_ZbHoI3gQTv4mR0CDLNQ,DQNYSjRH7AYthVmJ7P9-T" {
		t.Error("m is not equal to what expected (it's a bit long), got:", string(m))
	}
}
