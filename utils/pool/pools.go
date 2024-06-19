package pool

import (
	"git.golaxy.org/tiny/utils/generic"
	"git.golaxy.org/tiny/utils/types"
	"reflect"
)

var pools = map[string]*Pool{}

func Declare(i any) *Pool {
	return DeclareRT(reflect.TypeOf(i))
}

func DeclareRT(t reflect.Type) *Pool {
	pool := NewPoolRT(t)
	pools[pool.Name()] = pool
	return pool
}

func DeclareT[T any]() *Pool {
	pool := NewPoolT[T]()
	pools[pool.Name()] = pool
	return pool
}

func Get[T any](name string) *T {
	return pools[name].Get().(*T)
}

func Put[T any](name string, v *T) {
	*v = types.ZeroT[T]()
	pools[name].Put(v)
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
			TotalAlloc: pool.TotalAlloc(),
			TotalGet:   pool.TotalGet(),
			TotalPut:   pool.TotalPut(),
		})
	}

	return infos.Values()
}
