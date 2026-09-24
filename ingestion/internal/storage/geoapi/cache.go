package geoapi

import "sync"

type TypeID string

type TypesCache struct {
	types map[string]TypeID
	mu    sync.RWMutex
}

func NewTypesCache() *TypesCache {
	return &TypesCache{
		types: make(map[string]TypeID),
	}
}

func (c *TypesCache) Set(dataset string, id TypeID) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.types[dataset] = id
}

func (c *TypesCache) Get(dataset string) (TypeID, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.types[dataset]
	return id, ok
}
