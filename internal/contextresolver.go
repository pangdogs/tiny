package internal

import (
	"kit.golaxy.org/tiny/util"
)

// ContextResolver 上下文获取器
type ContextResolver interface {
	// ResolveContext 解析上下文
	ResolveContext() util.IfaceCache
}
