package tiny

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/runtime"
	"kit.golaxy.org/tiny/util/container"
)

func (_runtime *RuntimeBehavior) gc() {
	runtime.UnsafeContext(_runtime.ctx).GetInnerGC().GC()
	localevent.UnsafeEvent(&_runtime.eventUpdate).GC()
	localevent.UnsafeEvent(&_runtime.eventLateUpdate).GC()
}

// CollectGC 收集GC
func (_runtime *RuntimeBehavior) CollectGC(gc container.GC) {
}
