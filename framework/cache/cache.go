package cache

import (
	"gl.king.im/king-lib/framework/cache/store"
)

type Cache struct {
	store store.StoreInterface
}

func (c *Cache) Set(key, value string, secondsTtl int) error {
	return c.store.Set(key, value, secondsTtl)
}

func (c *Cache) Get(key string, defaultVal interface{}) (interface{}, error) {
	return c.store.Get(key, defaultVal)
}

func (c *Cache) Delete(key string) error {
	return c.store.Delete(key)
}

func (c *Cache) Has(key string) (bool, error) {
	return c.store.Has(key)
}

func New(store store.StoreInterface) (*Cache, error) {
	cacheIns := &Cache{}

	cacheIns.store = store

	return cacheIns, nil
}
