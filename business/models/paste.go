package models

import "time"

type Item struct {
	ID        string    `db:"id"`
	Paste     string    `db:"content" redis:"content"`
	CreatedAt time.Time `db:"created"`
}

type Items []Items
