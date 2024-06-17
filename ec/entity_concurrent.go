package ec

import (
	"fmt"
	"git.golaxy.org/tiny/internal/gctx"
	"git.golaxy.org/tiny/utils/uid"
)

// ConcurrentEntity 多线程安全的实体接口
type ConcurrentEntity interface {
	gctx.ConcurrentContextProvider
	fmt.Stringer

	// GetId 获取实体Id
	GetId() uid.Id
	// GetPrototype 获取实体原型
	GetPrototype() string
}
