package pool

import "git.golaxy.org/tiny/utils/types"

type PoolObject struct {
	Pool   *Pool
	Object any
}

type ManagedPoolObject interface {
	ManagedPoolObject(po PoolObject)
}

func ManagedGet[T any](managed ManagedPoolObject, pool *Pool) *T {
	if managed == nil {
		return types.NewT[T]()
	}

	obj := pool.Get().(*T)

	managed.ManagedPoolObject(PoolObject{
		Pool:   pool,
		Object: obj,
	})

	return obj
}
