package models

import "time"

type Item struct {
	ID        string    `db:"id"`
	Paste     string    `db:"content"`
	CreatedAt time.Time `db:"created"`
	Hash      string    `db:"hash"`
	IP        string    `db:"ip"`
	User      string    `db:"user"`
}

type Items []Items
