package repository

import (
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"polarite/resources"

	"github.com/dgraph-io/badger/v3"
)

func (d *Dependency) InsertPaste(ctx context.Context, paste Item) (Item, error) {
	compressedPaste, err := resources.CompressContent(paste.Paste)
	if err != nil {
		return Item{}, err
	}

	var existingItem Item
	var encodedMetadata bytes.Buffer
	err = gob.NewEncoder(&encodedMetadata).Encode(paste.Metadata)
	if err != nil {
		return Item{}, fmt.Errorf("encoding metadata to gob: %w", err)
	}

	err = d.DB.View(func(txn *badger.Txn) error {
		defer txn.Discard()

		existingItemId, err := txn.Get([]byte(paste.Hash))
		if err != nil && !errors.Is(err, badger.ErrKeyNotFound) {
			return fmt.Errorf("getting paste id by hash: %w", err)
		}

		if err == nil && existingItemId != nil {
			// Key already found
			existingItem.Hash = paste.Hash

			err := existingItemId.Value(func(val []byte) error {
				existingItem.ID = string(val)
				return nil
			})
			if err != nil {
				return fmt.Errorf("getting value for existing item id: %w", err)
			}

			return nil
		}

		return txn.Commit()
	})
	if err != nil {
		return Item{}, fmt.Errorf("getting existing data: %w", err)
	}

	if existingItem.ID != "" {
		return Item{
			ID:    existingItem.ID,
			Paste: paste.Paste,
			Hash:  paste.Hash,
		}, nil
	}

	batchWriter := d.DB.NewWriteBatch()

	err = batchWriter.SetEntry(&badger.Entry{
		Key:       []byte(paste.Hash),
		Value:     []byte(paste.ID),
		ExpiresAt: uint64(paste.Metadata.ExpiresAt.Unix()),
	})
	if err != nil {
		batchWriter.Cancel()
		return Item{}, fmt.Errorf("setting hash to id: %w", err)
	}

	err = batchWriter.SetEntry(&badger.Entry{
		Key:       []byte(paste.ID),
		Value:     compressedPaste,
		ExpiresAt: uint64(paste.Metadata.ExpiresAt.Unix()),
	})
	if err != nil {
		batchWriter.Cancel()
		return Item{}, fmt.Errorf("setting id to paste: %w", err)
	}

	err = batchWriter.SetEntry(&badger.Entry{
		Key:       []byte("metadata:" + paste.ID),
		Value:     encodedMetadata.Bytes(),
		ExpiresAt: uint64(paste.Metadata.ExpiresAt.Unix()),
	})
	if err != nil {
		batchWriter.Cancel()
		return Item{}, fmt.Errorf("setting id metadata: %w", err)
	}

	err = batchWriter.Flush()
	if err != nil {
		return Item{}, fmt.Errorf("inserting data: %w", err)
	}

	return paste, nil
}
