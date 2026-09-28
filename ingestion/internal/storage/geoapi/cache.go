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

type BusinessTypesCache struct {
	types map[string]string
	ids   map[string]string
	mu    sync.RWMutex
}

func NewBusinessTypesCache() *BusinessTypesCache {
	return &BusinessTypesCache{
		types: make(map[string]string),
		ids:   make(map[string]string),
	}
}

func (c *BusinessTypesCache) Set(id, slug string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.types[id] = slug
	c.ids[slug] = id
}

func (c *BusinessTypesCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.types)
}

func (c *BusinessTypesCache) ID(slug string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.ids[slug]
	return id, ok
}

func (c *BusinessTypesCache) Get(id string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	slug, ok := c.types[id]
	return slug, ok
}
