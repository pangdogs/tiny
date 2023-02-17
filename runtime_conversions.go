package tiny

import "kit.golaxy.org/tiny/util"

// GetInheritor 获取运行时的继承者
func GetInheritor[T any](runtime Runtime) T {
	return util.Cache2Iface[T](runtime.getOptions().Inheritor.Cache)
}
