package pool

import (
	"git.golaxy.org/tiny/utils/generic"
)

var pools = map[string]*Pool{}

func Declare[T any]() *Pool {
	pool := NewPool[T]()

	existed, ok := pools[pool.Name()]
	if ok {
		return existed
	}
	pools[pool.Name()] = pool

	return pool
}

func Get[T any](name string) *T {
	return pools[name].Get().(*T)
}

func Put[T any](name string, obj *T) {
	pools[name].Put(obj)
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
	TotalAlloc int64  `json:"total_alloc"`
	TotalGet   int64  `json:"total_get"`
	TotalPut   int64  `json:"total_put"`
}

func Stat() []Info {
	infos := generic.MakeSliceMap[string, Info]()

	for name, pool := range pools {
		infos.Add(name, Info{
			Name:       pool.Name(),
			TotalAlloc: pool.TotalAlloc(),
			TotalGet:   pool.TotalGet(),
			TotalPut:   pool.TotalPut(),
		})
	}

	return infos.Values()
}
