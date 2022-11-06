package repository

import "time"

type Item struct {
	ID       string
	Paste    []byte
	Hash     string
	Metadata Metadata
}

type Metadata struct {
	CreatorIP string
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Items []Items
