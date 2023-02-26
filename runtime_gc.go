package tiny

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/runtime"
)

func (_runtime *RuntimeBehavior) gc() {
	runtime.UnsafeContext(_runtime.ctx).GC()
	localevent.UnsafeEvent(&_runtime.eventUpdate).GC()
	localevent.UnsafeEvent(&_runtime.eventLateUpdate).GC()

	if _runtime.opts.CustomGC != nil {
		_runtime.opts.CustomGC(_runtime.opts.Inheritor.Iface)
	}
}
