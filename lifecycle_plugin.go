package tiny

import (
	"git.golaxy.org/tiny/runtime"
)

// LifecyclePluginInit 运行时上的插件初始化回调，插件实现此接口即可使用
type LifecyclePluginInit interface {
	InitP(ctx runtime.Context)
}

// LifecyclePluginShut 运行时上的插件结束回调，插件实现此接口即可使用
type LifecyclePluginShut interface {
	ShutP(ctx runtime.Context)
}
