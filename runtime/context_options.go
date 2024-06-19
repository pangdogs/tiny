package runtime

import (
	"context"
	"git.golaxy.org/tiny/plugin"
	"git.golaxy.org/tiny/utils/generic"
	"git.golaxy.org/tiny/utils/iface"
	"git.golaxy.org/tiny/utils/option"
)

type (
	RunningHandler = generic.DelegateAction2[Context, RunningState] // 运行状态变化处理器
)

// ContextOptions 创建运行时上下文的所有选项
type ContextOptions struct {
	CompositeFace     iface.Face[Context] // 扩展者，在扩展运行时上下文自身能力时使用
	Context           context.Context     // 父Context
	AutoRecover       bool                // 是否开启panic时自动恢复
	ReportError       chan error          // panic时错误写入的error channel
	PluginBundle      plugin.PluginBundle // 插件包
	UseObjectPool     bool                // 使用对象池
	UseObjectPoolSize int                 // 托管对象池初始大小
	RunningHandler    RunningHandler      // 运行状态变化处理器
}

type _ContextOption struct{}

// Default 默认值
func (_ContextOption) Default() option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		With.Context.CompositeFace(iface.Face[Context]{})(o)
		With.Context.Context(nil)(o)
		With.Context.PanicHandling(false, nil)(o)
		With.Context.PluginBundle(plugin.NewPluginBundle())(o)
		With.Context.UseObjectPool(false, 0)(o)
		With.Context.RunningHandler(nil)(o)
	}
}

// CompositeFace 扩展者，在扩展运行时上下文自身能力时使用
func (_ContextOption) CompositeFace(face iface.Face[Context]) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (_ContextOption) Context(ctx context.Context) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.Context = ctx
	}
}

// PanicHandling panic时的处理方式
func (_ContextOption) PanicHandling(autoRecover bool, reportError chan error) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.AutoRecover = autoRecover
		o.ReportError = reportError
	}
}

// PluginBundle 插件包
func (_ContextOption) PluginBundle(bundle plugin.PluginBundle) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.PluginBundle = bundle
	}
}

// UseObjectPool 使用对象池
func (_ContextOption) UseObjectPool(b bool, size int) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.UseObjectPool = b
		o.UseObjectPoolSize = size
	}
}

// RunningHandler 运行状态变化处理器
func (_ContextOption) RunningHandler(handler RunningHandler) option.Setting[ContextOptions] {
	return func(o *ContextOptions) {
		o.RunningHandler = handler
	}
}
