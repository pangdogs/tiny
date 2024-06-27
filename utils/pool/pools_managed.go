package pool

import (
	"git.golaxy.org/tiny/utils/types"
)

type PooledChunk struct {
	Pool  *Pool
	Chunk IChunk
}

type ManagedPooledChunk interface {
	ManagedGet(poolId uint32) any
	ManagedPooledChunk(pc PooledChunk)
}

func ManagedGet[T any](managed ManagedPooledChunk, pool *Pool) *T {
	if managed == nil {
		return types.NewT[T]()
	}

	obj := managed.ManagedGet(pool.Id())
	if obj == nil {
		managed.ManagedPooledChunk(PooledChunk{
			Pool:  pool,
			Chunk: pool.Get(),
		})
		obj = managed.ManagedGet(pool.Id())
		if obj == nil {
			panic("managed get pooled object failed")
		}
	}

	return obj.(*T)
}
