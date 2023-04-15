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

	return util.Cache2Iface[Context](ctxResolver.ResolveContext())
}
