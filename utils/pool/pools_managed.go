package pool

type PoolObject struct {
	Pool   *Pool
	Object any
}

type ManagedPoolObject interface {
	ManagedPoolObject(po PoolObject)
}

func ManagedGet[T any](managed ManagedPoolObject, pool *Pool) *T {
	obj := pool.Get().(*T)
	managed.ManagedPoolObject(PoolObject{
		Pool:   pool,
		Object: obj,
	})
	return obj
}
