package lxlib

import (
	"plugin"
)

type (
	Store struct {
		store *plugin.Plugin
	}
)

func NewStore(store *plugin.Plugin) (*Store, error) {
	if _, err := store.Lookup("Get"); err != nil {
		return nil, err
	}
	if _, err := store.Lookup("Set"); err != nil {
		return nil, err
	}
	return &Store {
		store: store,
	}, nil
}

func (this *Store) Get(key string) interface{} {
	fn, _ := this.store.Lookup("Get")
	return fn.(func(string) interface{})(key)
}

func (this *Store)  Set(key string, value interface{}) {
	fn, _ := this.store.Lookup("Set")
	fn.(func(string, interface{}))(key, value)
}