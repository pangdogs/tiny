package pool

import "git.golaxy.org/tiny/utils/types"

func NewBatch[T any](pool *Pool) *Batch[T] {
	batch := &Batch[T]{
		pool: pool,
	}
	return batch
}

func MakeBatch[T any](pool *Pool) Batch[T] {
	batch := Batch[T]{
		pool: pool,
	}
	return batch
}

type Batch[T any] struct {
	pool  *Pool
	chain []*T
}

func (c *Batch[T]) Cleanup() {
	for i := range c.chain {
		v := c.chain[i]
		*v = types.ZeroT[T]()
		c.pool.Put(v)
	}
	c.chain = nil
}

func (c *Batch[T]) Fetch() *T {
	v := c.pool.Get().(*T)
	c.chain = append(c.chain, v)
	return v
}
