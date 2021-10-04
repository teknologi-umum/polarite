// Insert a data to the database

package controllers

import (
	"context"
	"polarite/business/models"
	"time"

	"github.com/aidarkhanov/nanoid/v2"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4/pgxpool"
)

func InsertPasteToDB(db *pgxpool.Conn, body []byte) (models.Item, error) {
	defer db.Release()

	id, err := nanoid.New()
	if err != nil {
		return models.Item{}, err
	}

	creationTime := time.Now().Format(time.RFC3339)
	r, err := db.Query(context.Background(), "INSERT INTO paste (id, content, created) VALUES ($1, $2, $3) RETURNING id", id, string(body), creationTime)
	if err != nil {
		return models.Item{}, err
	}
	defer r.Close()

	var result models.Item
	err = pgxscan.ScanOne(&result, r)
	if err != nil {
		return models.Item{}, err
	}

	return result, nil
}
