package migration

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func Migrate(conn *pgxpool.Conn) error {
	defer conn.Release()

	r, err := conn.Query(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS paste (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			content TEXT NOT NULL,
			created TIMESTAMP NOT NULL
		);`,
	)
	if err != nil {
		return err
	}

	defer r.Close()
	return nil
}
