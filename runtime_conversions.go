package tiny

import "kit.golaxy.org/tiny/util"

// GetComposite 获取运行时的扩展者
func GetComposite[T any](runtime Runtime) T {
	return util.Cache2Iface[T](runtime.getOptions().CompositeFace.Cache)
}
