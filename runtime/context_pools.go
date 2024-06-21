package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/utils/pool"
)

// ManagedPoolObject 托管对象池，在运行时停止时自动解释放
func (ctx *ContextBehavior) ManagedPoolObject(po pool.PoolObject) {
	if !ctx.opts.UsePool {
		panic(fmt.Errorf("%w: not use object pool", ErrContext))
	}
	ctx.managedPoolObjects = append(ctx.managedPoolObjects, po)
}

// AutoUsePool 自动判断使用托管对象池
func (ctx *ContextBehavior) AutoUsePool() pool.ManagedPoolObject {
	if !ctx.opts.UsePool {
		return nil
	}
	return ctx.opts.CompositeFace.Iface
}

func (ctx *ContextBehavior) cleanManagedPoolObjects() {
	if !ctx.opts.UsePool {
		return
	}

	managedPoolObjects := ctx.managedPoolObjects
	ctx.managedPoolObjects = nil

	go func() {
		for i := range managedPoolObjects {
			po := &managedPoolObjects[i]
			po.Pool.Put(po.Object)
		}
	}()
}
