package runtime

import (
	"context"
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// ContextOptions 创建运行时上下文的所有选项
type ContextOptions struct {
	CompositeFace      util.Face[Context]                   // 扩展者，需要扩展运行时上下文自身能力时需要使用
	Context            context.Context                      // 父Context
	AutoRecover        bool                                 // 是否开启panic时自动恢复
	ReportError        chan error                           // panic时错误写入的error channel
	PersistIDGenerator ec.ID                                // 持久化ID生成器初始值
	StartedCallback    func(runtimeCtx Context)             // 启动运行时回调函数
	StoppingCallback   func(runtimeCtx Context)             // 开始停止运行时回调函数
	StoppedCallback    func(runtimeCtx Context)             // 完全停止运行时回调函数
	FrameBeginCallback func(runtimeCtx Context)             // 帧开始时的回调函数
	FrameEndCallback   func(runtimeCtx Context)             // 帧结束时的回调函数
	FaceAnyAllocator   container.Allocator[util.FaceAny]    // 自定义FaceAny内存分配器，用于提高性能
	HookAllocator      container.Allocator[localevent.Hook] // 自定义Hook内存分配器，用于提高性能
}

// ContextOption 创建运行时上下文的选项设置器
type ContextOption func(o *ContextOptions)

// WithContextOption 创建运行时上下文的所有选项设置器
type WithContextOption struct{}

// Default 默认值
func (WithContextOption) Default() ContextOption {
	return func(o *ContextOptions) {
		WithContextOption{}.CompositeFace(util.Face[Context]{})(o)
		WithContextOption{}.Context(nil)(o)
		WithContextOption{}.AutoRecover(false)(o)
		WithContextOption{}.ReportError(nil)(o)
		WithContextOption{}.PersistIDGenerator(0)(o)
		WithContextOption{}.StartedCallback(nil)(o)
		WithContextOption{}.StoppingCallback(nil)(o)
		WithContextOption{}.StoppedCallback(nil)(o)
		WithContextOption{}.FrameBeginCallback(nil)(o)
		WithContextOption{}.FrameEndCallback(nil)(o)
		WithContextOption{}.FaceAnyAllocator(container.DefaultAllocator[util.FaceAny]())(o)
		WithContextOption{}.HookAllocator(container.DefaultAllocator[localevent.Hook]())(o)
	}
}

// CompositeFace 扩展者，需要扩展运行时上下文自身功能时需要使用
func (WithContextOption) CompositeFace(v util.Face[Context]) ContextOption {
	return func(o *ContextOptions) {
		o.CompositeFace = v
	}
}

// Context 父Context
func (WithContextOption) Context(v context.Context) ContextOption {
	return func(o *ContextOptions) {
		o.Context = v
	}
}

// AutoRecover 是否开启panic时自动恢复
func (WithContextOption) AutoRecover(v bool) ContextOption {
	return func(o *ContextOptions) {
		o.AutoRecover = v
	}
}

// ReportError panic时错误写入的error channel
func (WithContextOption) ReportError(v chan error) ContextOption {
	return func(o *ContextOptions) {
		o.ReportError = v
	}
}

// PersistIDGenerator 持久化ID生成器初始值
func (WithContextOption) PersistIDGenerator(v ec.ID) ContextOption {
	return func(o *ContextOptions) {
		o.PersistIDGenerator = v
	}
}

// StartedCallback 启动运行时回调函数
func (WithContextOption) StartedCallback(v func(runtimeCtx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StartedCallback = v
	}
}

// StoppingCallback 开始停止运行时回调函数
func (WithContextOption) StoppingCallback(v func(runtimeCtx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppingCallback = v
	}
}

// StoppedCallback 完全停止运行时回调函数
func (WithContextOption) StoppedCallback(v func(runtimeCtx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppedCallback = v
	}
}

// FrameBeginCallback 帧更新开始时的回调函数
func (WithContextOption) FrameBeginCallback(v func(runtimeCtx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameBeginCallback = v
	}
}

// FrameEndCallback 帧更新结束时的回调函数
func (WithContextOption) FrameEndCallback(v func(runtimeCtx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameEndCallback = v
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能
func (WithContextOption) FaceAnyAllocator(v container.Allocator[util.FaceAny]) ContextOption {
	return func(o *ContextOptions) {
		if v == nil {
			panic("nil allocator")
		}
		o.FaceAnyAllocator = v
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能
func (WithContextOption) HookAllocator(v container.Allocator[localevent.Hook]) ContextOption {
	return func(o *ContextOptions) {
		if v == nil {
			panic("nil allocator")
		}
		o.HookAllocator = v
	}
}
