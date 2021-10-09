package controllers

import (
	"polarite/business/models"
)

type PasteController interface {
	ReadItemFromCache(id string) (models.Item, error)
	ReadItemFromDB(id string) (models.Item, error)
	ReadIDFromDB() ([]models.Item, error)
	ReadIDFromMemory() ([]string, error)
	InsertPasteToDB(body []byte) (models.Item, error)
	InsertPasteToCache(paste models.Item) error
	UpdateIDListFromDB(pastes []models.Item) ([]string, error)
	UpdateIDListFromCache(pastes []string, new string) (int, error)
}
