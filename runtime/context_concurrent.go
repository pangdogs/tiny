package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/internal/gctx"
	"git.golaxy.org/tiny/utils/async"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/iface"
)

// ConcurrentContextProvider 多线程安全的上下文提供者
type ConcurrentContextProvider = gctx.ConcurrentContextProvider

// ConcurrentContext 多线程安全的运行时上下文接口
type ConcurrentContext interface {
	gctx.ConcurrentContextProvider
	gctx.Context
	async.Caller
}

// Concurrent 获取多线程安全的运行时上下文
func Concurrent(provider gctx.ConcurrentContextProvider) ConcurrentContext {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetConcurrentContext())
}

func getCaller(provider gctx.ConcurrentContextProvider) async.Caller {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetConcurrentContext())
}
