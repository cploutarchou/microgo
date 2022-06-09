package cache

import (
	"github.com/dgraph-io/badger/v3"
	"time"
)

type BadgerCache struct {
	Connection *badger.DB
	Prefix     string
}

//Exists :  Check if the key exists
func (b *BadgerCache) Exists(str string) (bool, error) {
	_, err := b.Get(str)
	if err != nil {
		return false, err
	}
	return true, nil
}

//Get : Return Key values from Redis if it exists.
func (b *BadgerCache) Get(str string) (interface{}, error) {
	var fromCache []byte
	err := b.Connection.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(str))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			fromCache = append([]byte{}, val...)
			return nil
		})
		if err != nil {
			return nil
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	decoded, err := decode(string(fromCache))
	if err != nil {
		return nil, err
	}
	return decoded[str], nil
}

//Set : Set a key value in Redis
func (b *BadgerCache) Set(str string, value interface{}, expires ...int) error {
	entry := Entry{}
	entry[str] = value
	encoded, err := encode(entry)
	if err != nil {
		return err
	}
	if len(expires) > 0 {
		err = b.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded).WithTTL(time.Second * time.Duration(expires[0]))
			err = txn.SetEntry(e)
			return err
		})
	} else {
		err = b.Connection.Update(func(txn *badger.Txn) error {
			e := badger.NewEntry([]byte(str), encoded)
			err = txn.SetEntry(e)
			return err
		})
	}
	return nil
}

//Delete : Delete a key value in Redis.
func (b *BadgerCache) Delete(str string) error {
	err := b.Connection.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(str))
		return err
	})

	return err
}

//DeleteIfMatch : Delete a key values where match the with the key value
func (b *BadgerCache) DeleteIfMatch(str string) error {
	return b.deleteIfMatch(str)
}

//Clean : Delete all entries from redis.
func (b *BadgerCache) Clean() error {
	return b.deleteIfMatch("")
}
func (b *BadgerCache) deleteIfMatch(str string) error {

	deleteKeys := func(keysForDelete [][]byte) error {
		if err := b.Connection.Update(func(txn *badger.Txn) error {
			for _, key := range keysForDelete {
				if err := txn.Delete(key); err != nil {
					return err
				}
			}
			return nil
		}); err != nil {
			return err
		}
		return nil
	}

	collectSize := 100000

	err := b.Connection.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.AllVersions = false
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)

		keysForDelete := make([][]byte, 0, collectSize)
		keysCollected := 0

		for it.Seek([]byte(str)); it.ValidForPrefix([]byte(str)); it.Next() {
			key := it.Item().KeyCopy(nil)
			keysForDelete = append(keysForDelete, key)
			keysCollected++
			if keysCollected == collectSize {
				if err := deleteKeys(keysForDelete); err != nil {
					return err
				}
			}
		}

		if keysCollected > 0 {
			if err := deleteKeys(keysForDelete); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}
