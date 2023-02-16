package runtime

import "kit.golaxy.org/tiny/util/container"

type _InnerGC interface {
	getInnerGC() container.GC
}
