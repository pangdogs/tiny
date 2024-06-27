package pool

import (
	"git.golaxy.org/tiny/utils/types"
	"reflect"
	"sync"
	"sync/atomic"
)

func NewPool[T any]() *Pool {
	pool := &Pool{}
	pool.reflectType = reflect.TypeFor[T]()
	pool.name = types.FullNameRT(pool.reflectType)
	pool.pool.New = func() any {
		atomic.AddInt64(&pool.allocNum, 1)
		return types.NewT[T]()
	}
	pool.zero = func(obj any) any {
		*(obj.(*T)) = types.ZeroT[T]()
		return obj
	}
	return pool
}

type Pool struct {
	name                     string
	pool                     sync.Pool
	reflectType              reflect.Type
	zero                     func(obj any) any
	allocNum, getNum, putNum int64
}

func (p *Pool) Name() string {
	return p.name
}

func (p *Pool) ReflectType() reflect.Type {
	return p.reflectType
}

func (p *Pool) Put(obj any) {
	p.pool.Put(p.zero(obj))
	atomic.AddInt64(&p.putNum, 1)
}

func (p *Pool) Get() any {
	v := p.pool.Get()
	atomic.AddInt64(&p.getNum, 1)
	return v
}

func (p *Pool) TotalAlloc() int64 {
	return atomic.LoadInt64(&p.allocNum)
}

func (p *Pool) TotalGet() int64 {
	return atomic.LoadInt64(&p.getNum)
}

func (p *Pool) TotalPut() int64 {
	return atomic.LoadInt64(&p.putNum)
}
