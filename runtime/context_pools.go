package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/utils/pool"
)

// ManagedGet 从托管对象池中获取对象
func (ctx *ContextBehavior) ManagedGet(poolId uint32) any {
	if !ctx.opts.UsePool {
		panic(fmt.Errorf("%w: not use object pool", ErrContext))
	}
	pc, ok := ctx.managedPooledUsed[poolId]
	if !ok {
		return nil
	}
	return pc.Chunk.Get()
}

// ManagedPooledChunk 托管对象池，在运行时停止时自动解释放
func (ctx *ContextBehavior) ManagedPooledChunk(pc pool.PooledChunk) {
	if !ctx.opts.UsePool {
		panic(fmt.Errorf("%w: not use object pool", ErrContext))
	}
	idx := len(ctx.managedPooledChunk)
	ctx.managedPooledChunk = append(ctx.managedPooledChunk, pc)
	ctx.managedPooledUsed[pc.Pool.Id()] = &ctx.managedPooledChunk[idx]
}

// AutoUsePool 自动判断使用托管对象池
func (ctx *ContextBehavior) AutoUsePool() pool.ManagedPooledChunk {
	if !ctx.opts.UsePool {
		return nil
	}
	return ctx.opts.CompositeFace.Iface
}

func (ctx *ContextBehavior) cleanManagedPoolObjects() {
	if !ctx.opts.UsePool {
		return
	}

	managedPoolObjects := ctx.managedPooledChunk
	ctx.managedPooledChunk = nil
	ctx.managedPooledUsed = nil

	go func() {
		for i := range managedPoolObjects {
			pc := &managedPoolObjects[i]
			pc.Pool.Put(pc.Chunk)
		}
	}()
}
