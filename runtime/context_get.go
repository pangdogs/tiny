package runtime

import (
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/util"
)

// Get 获取运行时上下文
func Get(ctxResolver ec.ContextResolver) Context {
	if ctxResolver == nil {
		panic("nil ctxResolver")
	}

	ctx := ec.UnsafeContextResolver(ctxResolver).GetContext()
	if ctx == util.NilIfaceCache {
		panic("nil context")
	}

	return util.Cache2Iface[Context](ctx)
}

// TryGet 尝试获取运行时上下文
func TryGet(ctxResolver ec.ContextResolver) (Context, bool) {
	if ctxResolver == nil {
		return nil, false
	}

	ctx := ec.UnsafeContextResolver(ctxResolver).GetContext()
	if ctx == util.NilIfaceCache {
		return nil, false
	}

	return util.Cache2Iface[Context](ctx), true
}
