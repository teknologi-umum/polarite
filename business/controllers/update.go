package controllers

import (
	"polarite/business/models"
	"strings"

	"github.com/allegro/bigcache/v3"
)

func UpdateIDListFromDB(memory *bigcache.BigCache, pastes []models.Item) ([]string, error) {
	var temp []string
	for i := 0; i < len(pastes); i++ {
		temp = append(temp, pastes[i].ID)
	}

	s := strings.Join(temp, ",")
	err := memory.Set("ids", []byte(s))
	if err != nil {
		return []string{}, err
	}

	return temp, nil
}

func UpdateIDListFromCache(memory *bigcache.BigCache, pastes []string, new string) (int, error) {
	pastes = append(pastes, new)
	s := strings.Join(pastes, ",")
	err := memory.Set("ids", []byte(s))
	if err != nil {
		return 0, err
	}

	return len(pastes), nil
}
