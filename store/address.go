package store

import (
	"encoding/json"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type Address struct {
	ID      string
	UserID  string
	Line1   string
	City    string
	ZipCode string
}

const addressPrefix = "addr:"

func AddAddressToUser(userID, line1, city, zip string) (*Address, error) {
	addr := &Address{
		ID:      uuid.NewString(),
		UserID:  userID,
		Line1:   line1,
		City:    city,
		ZipCode: zip,
	}
	data, err := json.Marshal(addr)
	if err != nil {
		return nil, err
	}
	err = GetDB().Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(addressPrefix+addr.ID), data)
	})
	return addr, err
}

func ListAddressesByUserID(userID string) ([]*Address, error) {
	addresses := []*Address{}
	err := GetDB().View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(addressPrefix)); it.ValidForPrefix([]byte(addressPrefix)); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var a Address
				if err := json.Unmarshal(v, &a); err == nil && a.UserID == userID {
					addresses = append(addresses, &a)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return addresses, err
}
