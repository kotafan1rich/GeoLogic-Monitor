package twogis

import (
	"sync"
	"time"
)

type CheckpointsCache struct {
	checkpoints map[string]time.Time
	mu          sync.RWMutex
}

func NewCheckpointsCache() *CheckpointsCache {
	return &CheckpointsCache{
		checkpoints: make(map[string]time.Time),
	}
}

func (c *CheckpointsCache) Set(businessType string, date time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.checkpoints[businessType] = date
}

func (c *CheckpointsCache) Get(dataset string) (time.Time, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	id, ok := c.checkpoints[dataset]
	return id, ok
}
