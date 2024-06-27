package pool

import (
	"fmt"
	"git.golaxy.org/tiny/utils/generic"
)

var (
	poolNames = map[string]*Pool{}
	poolIds   = map[uint32]*Pool{}
)

func Declare[T any](chunkSize int64) *Pool {
	pool := NewPool[T](chunkSize)

	existed, ok := poolNames[pool.Name()]
	if ok {
		if existed.chunkSize < chunkSize {
			existed.chunkSize = chunkSize
		}
		return existed
	}

	existed, ok = poolIds[pool.Id()]
	if ok {
		panic(fmt.Errorf("pool(%d) has already been declared by %q", pool.Id(), existed.Name()))
	}

	poolNames[pool.Name()] = pool
	poolIds[pool.Id()] = pool

	return pool
}

func Get[T any](name string) *Chunk[T] {
	return poolNames[name].Get().(*Chunk[T])
}

func GetById[T any](id uint32) *Chunk[T] {
	return poolIds[id].Get().(*Chunk[T])
}

func Put[T any](name string, chunk *Chunk[T]) {
	poolNames[name].Put(chunk)
}

func PutById[T any](id uint32, chunk *Chunk[T]) {
	poolIds[id].Put(chunk)
}

func Range(fun generic.Func2[string, *Pool, bool]) {
	for name, pool := range poolNames {
		if !fun.Exec(name, pool) {
			return
		}
	}
}

func Each(fun generic.Action2[string, *Pool]) {
	for name, pool := range poolNames {
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

	for name, pool := range poolNames {
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
