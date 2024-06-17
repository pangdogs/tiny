package tiny

import (
	"git.golaxy.org/tiny/runtime"
)

func (rt *RuntimeBehavior) gc() {
	runtime.UnsafeContext(rt.ctx).GC()
	rt.opts.CustomGC.Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError(), nil, rt.opts.CompositeFace.Iface)
}
