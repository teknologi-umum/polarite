package controllers

import (
	"polarite/business/models"

	"github.com/jmoiron/sqlx"
)

type PasteController interface {
	ReadItemFromCache(id string) (models.Item, error)
	ReadItemFromDB(db *sqlx.Conn, id string) (models.Item, error)
	ReadIDFromDB(db *sqlx.Conn) ([]models.Item, error)
	ReadIDFromMemory() ([]string, error)
	InsertPasteToDB(db *sqlx.Conn, body []byte) (models.Item, error)
	InsertPasteToCache(paste models.Item) error
	UpdateIDListFromDB(pastes []models.Item) ([]string, error)
	UpdateIDListFromCache(pastes []string, new string) (int, error)
}
