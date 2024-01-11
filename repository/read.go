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

func (d *Dependency) GetItemById(ctx context.Context, id string) (Item, error) {
	_, span := tracer.Start(ctx, "GetItemById")
	defer span.End()

	var item = Item{ID: id}

	err := d.DB.View(func(txn *badger.Txn) error {
		compressedPaste, err := txn.Get([]byte(id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				txn.Discard()
				return err
			}

			txn.Discard()
			return fmt.Errorf("getting compressed paste: %w", err)
		}

		err = compressedPaste.Value(func(val []byte) error {
			decompressed, err := resources.DecompressContent(val)
			if err != nil {
				return fmt.Errorf("decompressing: %w", err)
			}

			item.Paste = decompressed
			return nil
		})
		if err != nil {
			txn.Discard()
			return err
		}

		encodedMetadata, err := txn.Get([]byte("metadata:" + id))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				txn.Discard()
				return err
			}

			txn.Discard()
			return fmt.Errorf("getting encoded metadata: %w", err)
		}

		err = encodedMetadata.Value(func(val []byte) error {
			var metadata Metadata
			err := gob.NewDecoder(bytes.NewReader(val)).Decode(&metadata)
			if err != nil {
				return fmt.Errorf("decoding: %w", err)
			}

			item.Metadata = metadata
			return nil
		})
		if err != nil {
			txn.Discard()
			return err
		}

		txn.Discard()
		return nil
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return Item{}, ErrNotFound
		}

		return Item{}, fmt.Errorf("reading database: %w", err)
	}

	return item, nil
}

func (d *Dependency) ReadHash(ctx context.Context, h string) (exists bool, id string, err error) {
	_, span := tracer.Start(ctx, "ReadHash")
	defer span.End()

	err = d.DB.View(func(txn *badger.Txn) error {
		defer txn.Discard()

		valueId, err := txn.Get([]byte(h))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				return err
			}

			return fmt.Errorf("getting compressed paste: %w", err)
		}

		err = valueId.Value(func(val []byte) error {
			id = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return false, "", nil
		}

		return false, "", fmt.Errorf("reading database: %w", err)
	}

	return true, id, nil
}
