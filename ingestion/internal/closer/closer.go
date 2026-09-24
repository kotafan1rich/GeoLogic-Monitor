package closer

import (
	"context"
	"errors"
	"sync"
)

type closeFn struct {
	name string
	fn   func(context.Context) error
}

type closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFn
}

var Closer = &closer{}

func Add(name string, fn func(context.Context) error) {
	Closer.add(name, fn)
}

func CloseAll(ctx context.Context) error {
	return Closer.closeAll(ctx)
}

func (c *closer) add(name string, fn func(context.Context) error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.funcs = append(c.funcs, closeFn{name: name, fn: fn})
}

func (c *closer) closeAll(ctx context.Context) error {
	var res error

	c.once.Do(func() {
		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		if len(funcs) == 0 {
			return
		}

		var errs []error

		for i := len(funcs) - 1; i >= 0; i-- {
			if err := funcs[i].fn(ctx); err != nil {
				errs = append(errs, err)
			}
		}

		res = errors.Join(errs...)
	})

	return res
}
