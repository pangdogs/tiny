package pool

import (
	"context"
	"git.golaxy.org/tiny/utils/types"
	common_pool "github.com/jolestar/go-commons-pool/v2"
	"sync/atomic"
)

func NewPool[T any]() *Pool {
	pool := &Pool{}
	pool.name = types.FullNameT[T]()

	pool.pool = common_pool.NewObjectPool(
		context.Background(),
		common_pool.NewPooledObjectFactory(
			func(ctx context.Context) (interface{}, error) {
				atomic.AddInt64(&pool.allocNum, 1)
				return types.NewT[T](), nil
			},
			nil,
			nil,
			nil,
			func(ctx context.Context, obj *common_pool.PooledObject) error {
				*(obj.Object.(*T)) = types.ZeroT[T]()
				atomic.AddInt64(&pool.putNum, 1)
				return nil
			},
		),
		&common_pool.ObjectPoolConfig{
			LIFO:                     false,
			MaxTotal:                 -1,
			MaxIdle:                  0,
			MinIdle:                  0,
			TestOnCreate:             false,
			TestOnBorrow:             false,
			TestOnReturn:             false,
			TestWhileIdle:            false,
			BlockWhenExhausted:       common_pool.DefaultBlockWhenExhausted,
			MinEvictableIdleTime:     common_pool.DefaultMinEvictableIdleTime,
			SoftMinEvictableIdleTime: common_pool.DefaultSoftMinEvictableIdleTime,
			NumTestsPerEvictionRun:   common_pool.DefaultNumTestsPerEvictionRun,
			EvictionPolicyName:       common_pool.DefaultEvictionPolicyName,
			TimeBetweenEvictionRuns:  common_pool.DefaultTimeBetweenEvictionRuns,
			EvictionContext:          context.Background(),
		},
	)

	return pool
}

type Pool struct {
	name                     string
	pool                     *common_pool.ObjectPool
	allocNum, getNum, putNum int64
}

func (p *Pool) Name() string {
	return p.name
}

func (p *Pool) Prepare(num int64) {
	p.pool.Config.MinIdle = int(num)
	p.pool.Config.MaxIdle = int(num)
	p.pool.PreparePool(context.Background())
}

func (p *Pool) Put(obj any) {
	p.pool.ReturnObject(context.Background(), obj)
}

func (p *Pool) Get() any {
	v, _ := p.pool.BorrowObject(context.Background())
	atomic.AddInt64(&p.allocNum, 1)
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
