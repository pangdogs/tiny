package ec

import (
	"kit.golaxy.org/tiny/util"
)

// GetInheritor 获取实体的继承者
func GetInheritor[T any](entity Entity) T {
	return util.Cache2Iface[T](entity.getOptions().Inheritor.Cache)
}
