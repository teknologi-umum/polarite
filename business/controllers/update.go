package controllers

import (
	"polarite/business/models"
	"strings"

	"github.com/allegro/bigcache/v3"
)

func UpdateIDList(memory *bigcache.BigCache, pastes []models.Item) (int, error) {
	var temp []string
	for i := 0; i < len(pastes); i++ {
		temp = append(temp, pastes[i].ID)
	}

	s := strings.Join(temp, ",")
	err := memory.Set("ids", []byte(s))
	if err != nil {
		return 0, err
	}

	return len(s), nil
}
