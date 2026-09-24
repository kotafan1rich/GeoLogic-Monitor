package errs

import (
	"errors"
	"fmt"
	"log/slog"
	"sync"
)

const DefaultSampleSize = 5

type Collector struct {
	mu     sync.Mutex
	limit  int
	total  int
	sample []error
}

func NewCollector(limit int) *Collector {
	if limit <= 0 {
		limit = DefaultSampleSize
	}

	return &Collector{limit: limit, sample: make([]error, 0, limit)}
}

func (c *Collector) Add(err error) {
	if err == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.total++

	if len(c.sample) < c.limit {
		c.sample = append(c.sample, err)
	}
}

func (c *Collector) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.total
}

func (c *Collector) Err() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.total == 0 {
		return nil
	}

	joined := make([]error, 0, len(c.sample)+1)
	joined = append(joined, c.sample...)

	if hidden := c.total - len(c.sample); hidden > 0 {
		joined = append(joined, fmt.Errorf("and %d more error(s)", hidden))
	}

	return errors.Join(joined...)
}

func (c *Collector) Attr() slog.Attr {
	if err := c.Err(); err != nil {
		return slog.Any("error", err)
	}

	return slog.Attr{}
}
