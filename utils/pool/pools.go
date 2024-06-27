package pool

import (
	"git.golaxy.org/tiny/utils/generic"
)

var pools = map[string]*Pool{}

func Declare[T any](chunkSize int64) *Pool {
	pool := NewPool[T](chunkSize)

	existed, ok := pools[pool.Name()]
	if ok {
		if existed.chunkSize < chunkSize {
			existed.chunkSize = chunkSize
		}
		return existed
	}
	pools[pool.Name()] = pool

	return pool
}

func Get[T any](name string) *Chunk[T] {
	return pools[name].Get().(*Chunk[T])
}

func Put[T any](name string, chunk *Chunk[T]) {
	pools[name].Put(chunk)
}

func Range(fun generic.Func2[string, *Pool, bool]) {
	for name, pool := range pools {
		if !fun.Exec(name, pool) {
			return
		}
	}
}

func Each(fun generic.Action2[string, *Pool]) {
	for name, pool := range pools {
		fun.Exec(name, pool)
	}
}

type Info struct {
	Name       string `json:"name"`
	ChunkSize  int64  `json:"chunk_size"`
	TotalAlloc int64  `json:"total_alloc"`
	TotalGet   int64  `json:"total_get"`
	TotalPut   int64  `json:"total_put"`
}

func Stat() []Info {
	infos := generic.MakeSliceMap[string, Info]()

	for name, pool := range pools {
		infos.Add(name, Info{
			Name:       pool.Name(),
			ChunkSize:  pool.ChunkSize(),
			TotalAlloc: pool.TotalAlloc(),
			TotalGet:   pool.TotalGet(),
			TotalPut:   pool.TotalPut(),
		})
	}

	return infos.Values()
}
