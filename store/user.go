package store

import (
	"encoding/json"
	"errors"

	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
)

type User struct {
	ID       string
	Name     string
	Email    string
	Password string
}

const (
	userPrefix     = "user:"
	emailIndexPref = "email:"
)

var (
	ErrEmailExists = errors.New("email already registered")
)

func CreateUser(name, email, password string) (*User, error) {
	if _, err := findUserIDByEmail(email); err == nil {
		return nil, ErrEmailExists
	}

	user := &User{
		ID:       uuid.NewString(),
		Name:     name,
		Email:    email,
		Password: password,
	}
	data, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	err = GetDB().Update(func(txn *badger.Txn) error {
		if err := txn.Set([]byte(userPrefix+user.ID), data); err != nil {
			return err
		}
		return txn.Set([]byte(emailIndexPref+email), []byte(user.ID))
	})

	return user, err
}

func FindUserByID(id string) (*User, error) {
	var user User
	err := GetDB().View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(userPrefix + id))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &user)
		})
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func findUserIDByEmail(email string) (string, error) {
	var id string
	err := GetDB().View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(emailIndexPref + email))
		if err != nil {
			return err
		}
		return item.Value(func(v []byte) error {
			id = string(v)
			return nil
		})
	})
	return id, err
}

func FindUserByEmail(email string) (*User, error) {
	id, err := findUserIDByEmail(email)
	if err != nil {
		return nil, err
	}
	return FindUserByID(id)
}

func ListUsers() ([]*User, error) {
	users := []*User{}
	err := GetDB().View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek([]byte(userPrefix)); it.ValidForPrefix([]byte(userPrefix)); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				var u User
				if err := json.Unmarshal(v, &u); err == nil {
					users = append(users, &u)
				}
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return users, err
}

func UpdateUser(u *User) error {
	data, err := json.Marshal(u)
	if err != nil {
		return err
	}
	return GetDB().Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(userPrefix+u.ID), data)
	})
}
