package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/internal/gctx"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/iface"
)

// CurrentContextProvider 当前上下文提供者
type CurrentContextProvider = gctx.CurrentContextProvider

// Current 获取当前运行时上下文
func Current(provider gctx.CurrentContextProvider) Context {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetCurrentContext())
}

func getRuntimeContext(provider gctx.CurrentContextProvider) Context {
	if provider == nil {
		panic(fmt.Errorf("%w: %w: provider is nil", ErrContext, exception.ErrArgs))
	}
	return iface.Cache2Iface[Context](provider.GetCurrentContext())
}
