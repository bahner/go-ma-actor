package db

import (
	"errors"
	"strings"

	badger "github.com/dgraph-io/badger/v3"
)

var (
	ErrKeyNotFound = errors.New("key not found")
)

// Returns the DID . Returns the input if the node does not exist
// This is used before we know in an Entity exists or not. It can be used anywhere.
func Get(key []byte) ([]byte, error) {

	var value []byte

	err := DB().View(func(txn *badger.Txn) error {
		item, err := txn.Get(key)
		if err != nil {
			return err // Key not found or other error
		}
		err = item.Value(func(val []byte) error {
			value = val
			return nil
		})
		return err
	})

	if err != nil {
		return nil, err
	}
	return value, nil
}

// Removes a node from the database if it exists. Must be a DID
// The prefix has a name, to make it obvious. eg. "entity:nick:" or "peer:nick:"
func Delete(key []byte) error {

	return DB().Update(func(txn *badger.Txn) error {
		return txn.Delete(key)
	})
}

// Sets a key to the value in the database
// Remember to use a prefix for the key. eg. "entity:nick:" or "peer:nick:"
func Set(key []byte, value []byte) error {

	return DB().Update(func(txn *badger.Txn) error {
		return txn.Set(key, value)
	})
}

// Sets a key to the value in the database, but
// it deletes the old key first if it exists.
func Upsert(prefix []byte, key []byte, value []byte) error {
	return DB().Update(func(txn *badger.Txn) error {
		oldKey, err := Lookup(value)
		if err == nil {
			if delErr := txn.Delete(oldKey); delErr != nil {
				return delErr
			}
		} else if err != ErrKeyNotFound && err != badger.ErrKeyNotFound {
			return err
		}
		fullKey := append(prefix, key...)
		return txn.Set(fullKey, value)
	})
}

// Try to find a key for a value. Returns the key first key found for the value.
func Lookup(value []byte) ([]byte, error) {
	var foundKey []byte

	err := DB().View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			err := item.Value(func(val []byte) error {
				if string(val) == string(value) {
					foundKey = item.Key()
					return nil
				}
				return nil
			})
			if err != nil {
				return err
			}
			if foundKey != nil {
				break // Exit the loop early if the key is found
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	if foundKey == nil {
		return nil, ErrKeyNotFound
	}
	return foundKey, nil
}

// Returns a map of all keys with a given prefix. It removes the prefix from the key.
func Keys(prefix string) (map[string]string, error) {
	nicks := make(map[string]string)

	err := DB().View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()

			err := item.Value(func(val []byte) error {
				nick := strings.TrimPrefix(string(key), prefix)
				nicks[nick] = string(val)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return nicks, nil
}
