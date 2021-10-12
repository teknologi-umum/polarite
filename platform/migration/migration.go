package migration

import (
	"context"

	"github.com/jmoiron/sqlx"
)

func Migrate(conn *sqlx.Conn) error {
	defer conn.Close()

	r, err := conn.QueryContext(
		context.Background(),
		`CREATE TABLE IF NOT EXISTS paste (
			id VARCHAR(36) PRIMARY KEY NOT NULL,
			content MEDIUMTEXT NOT NULL,
			created TIMESTAMP NOT NULL,
			hash VARCHAR(255) UNIQUE NOT NULL,
			ip VARCHAR(20) NOT NULL,
			user VARCHAR(255) NOT NULL
		);`,
	)
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}
