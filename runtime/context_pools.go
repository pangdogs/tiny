package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/utils/pool"
)

// ManagedPoolObject 托管对象池，在运行时停止时自动解释放
func (ctx *ContextBehavior) ManagedPoolObject(po pool.PoolObject) {
	if !ctx.opts.UseObjectPool {
		panic(fmt.Errorf("%w: not use object pool", ErrContext))
	}
	ctx.managedPoolObjects = append(ctx.managedPoolObjects, po)
}

// AutoManagedPoolObject 自动判断托管对象池
func (ctx *ContextBehavior) AutoManagedPoolObject() pool.ManagedPoolObject {
	if !ctx.opts.UseObjectPool {
		return nil
	}
	return ctx.opts.CompositeFace.Iface
}

func (ctx *ContextBehavior) cleanManagedPoolObjects() {
	if !ctx.opts.UseObjectPool {
		return
	}
	for i := range ctx.managedPoolObjects {
		po := &ctx.managedPoolObjects[i]
		po.Pool.Put(po.Object)
	}
	ctx.managedHooks = nil
}
