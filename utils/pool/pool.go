package pool

import (
	"git.golaxy.org/tiny/utils/types"
	"hash/fnv"
	"reflect"
	"sync"
	"sync/atomic"
)

func NewPool[T any](chunkSize int64) *Pool {
	pool := &Pool{}
	pool.reflectType = types.TypeFor[T]()
	pool.name = types.FullNameRT(pool.reflectType)
	pool.id = makePoolId(pool.reflectType)
	pool.chunkSize = chunkSize
	pool.pool.New = func() any {
		atomic.AddInt64(&pool.allocNum, 1)
		return &Chunk[T]{
			Objects: make([]T, pool.chunkSize),
			Pos:     0,
		}
	}
	pool.zero = func(chunk any) any {
		c := chunk.(*Chunk[T])
		clearSlice(c.Objects)
		c.Pos = 0
		return chunk
	}
	return pool
}

type IChunk interface {
	Get() any
}

type Chunk[T any] struct {
	Objects []T
	Pos     int32
}

func (c *Chunk[T]) Get() any {
	if c.Pos >= int32(len(c.Objects)) {
		return nil
	}
	obj := &c.Objects[c.Pos]
	c.Pos++
	return obj
}

type Pool struct {
	id                       uint32
	name                     string
	reflectType              reflect.Type
	chunkSize                int64
	pool                     sync.Pool
	zero                     func(any) any
	allocNum, getNum, putNum int64
}

func (p *Pool) Id() uint32 {
	return p.id
}

func (p *Pool) Name() string {
	return p.name
}

func (p *Pool) ReflectType() reflect.Type {
	return p.reflectType
}

func (p *Pool) ChunkSize() int64 {
	return p.chunkSize
}

func (p *Pool) Put(chunk IChunk) {
	p.pool.Put(p.zero(chunk))
	atomic.AddInt64(&p.putNum, 1)
}

func (p *Pool) Get() IChunk {
	chunk := p.pool.Get().(IChunk)
	atomic.AddInt64(&p.getNum, 1)
	return chunk
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

func makePoolId(rt reflect.Type) uint32 {
	hash := fnv.New32a()
	if rt.PkgPath() == "" || rt.Name() == "" {
		panic("unsupported type")
	}
	hash.Write([]byte(types.FullNameRT(rt)))
	return hash.Sum32()
}

func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}
