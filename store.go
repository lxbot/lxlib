package lxlib

import (
	"errors"
	"plugin"
)

type (
	Store struct {
		store *plugin.Plugin
	}
)

func NewStore(store *plugin.Plugin) (*Store, error) {
	fn, err := store.Lookup("Get");
	if err != nil {
		return nil, err
	}
	switch fn.(type) {
	case func(string) interface{}:
		// NOTE: NOP
		break
	default:
		return nil, errors.New("invalid store")
	}
	fn, err = store.Lookup("Get");
	if err != nil {
		return nil, err
	}
	switch fn.(type) {
	case func(string, interface{}):
		// NOTE: NOP
		break
	default:
		return nil, errors.New("invalid store")
	}

	return &Store{store}, nil
}

func (this *Store) Get(key string) interface{} {
	fn, _ := this.store.Lookup("Get")
	return fn.(func(string) interface{})(key)
}

func (this *Store)  Set(key string, value interface{}) {
	fn, _ := this.store.Lookup("Set")
	fn.(func(string, interface{}))(key, value)
}