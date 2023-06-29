package repository

import (
	"context"
	"errors"

	"github.com/dgraph-io/badger/v3"
)

type Dependency struct {
	DB *badger.DB
}

type Paste interface {
	GetItemById(ctx context.Context, id string) (Item, error)
	ReadHash(ctx context.Context, h string) (exists bool, id string, err error)
	InsertPaste(ctx context.Context, paste Item) (Item, error)
}

var ErrIDDuplicate = errors.New("generated id is duplicated, please try again")
var ErrNotFound = errors.New("not found")
var ErrNoID = errors.New("an ID needs to be supplied")
var ErrBodyTooBig = errors.New("body supplied exceeded the maximum size of 5 MB")
var ErrNoAuthHeader = errors.New("authorization headers must be supplied")
