package runtime

import (
	"context"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/uid"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// WithOption 所有选项设置器
type WithOption struct{}

// ContextOptions 创建运行时上下文的所有选项
type ContextOptions struct {
	CompositeFace      util.Face[Context]                   // 扩展者，需要扩展运行时上下文自身能力时需要使用
	Context            context.Context                      // 父Context
	AutoRecover        bool                                 // 是否开启panic时自动恢复
	ReportError        chan error                           // panic时错误写入的error channel
	PersistIdGenerator uid.Id                               // 持久化Id生成器初始值
	StartedCallback    func(ctx Context)                    // 启动运行时回调函数
	StoppingCallback   func(ctx Context)                    // 开始停止运行时回调函数
	StoppedCallback    func(ctx Context)                    // 完全停止运行时回调函数
	FrameBeginCallback func(ctx Context)                    // 帧开始时的回调函数
	FrameEndCallback   func(ctx Context)                    // 帧结束时的回调函数
	FaceAnyAllocator   container.Allocator[util.FaceAny]    // 自定义FaceAny内存分配器，用于提高性能
	HookAllocator      container.Allocator[localevent.Hook] // 自定义Hook内存分配器，用于提高性能
}

// ContextOption 创建运行时上下文的选项设置器
type ContextOption func(o *ContextOptions)

// Default 默认值
func (WithOption) Default() ContextOption {
	return func(o *ContextOptions) {
		WithOption{}.CompositeFace(util.Face[Context]{})(o)
		WithOption{}.Context(nil)(o)
		WithOption{}.AutoRecover(false)(o)
		WithOption{}.ReportError(nil)(o)
		WithOption{}.PersistIdGenerator(0)(o)
		WithOption{}.StartedCallback(nil)(o)
		WithOption{}.StoppingCallback(nil)(o)
		WithOption{}.StoppedCallback(nil)(o)
		WithOption{}.FrameBeginCallback(nil)(o)
		WithOption{}.FrameEndCallback(nil)(o)
		WithOption{}.FaceAnyAllocator(container.DefaultAllocator[util.FaceAny]())(o)
		WithOption{}.HookAllocator(container.DefaultAllocator[localevent.Hook]())(o)
	}
}

// CompositeFace 扩展者，需要扩展运行时上下文自身功能时需要使用
func (WithOption) CompositeFace(face util.Face[Context]) ContextOption {
	return func(o *ContextOptions) {
		o.CompositeFace = face
	}
}

// Context 父Context
func (WithOption) Context(ctx context.Context) ContextOption {
	return func(o *ContextOptions) {
		o.Context = ctx
	}
}

// AutoRecover 是否开启panic时自动恢复
func (WithOption) AutoRecover(b bool) ContextOption {
	return func(o *ContextOptions) {
		o.AutoRecover = b
	}
}

// ReportError panic时错误写入的error channel
func (WithOption) ReportError(ch chan error) ContextOption {
	return func(o *ContextOptions) {
		o.ReportError = ch
	}
}

// PersistIdGenerator 持久化Id生成器初始值
func (WithOption) PersistIdGenerator(v uid.Id) ContextOption {
	return func(o *ContextOptions) {
		o.PersistIdGenerator = v
	}
}

// StartedCallback 启动运行时回调函数
func (WithOption) StartedCallback(v func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StartedCallback = v
	}
}

// StoppingCallback 开始停止运行时回调函数
func (WithOption) StoppingCallback(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppingCallback = fn
	}
}

// StoppedCallback 完全停止运行时回调函数
func (WithOption) StoppedCallback(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppedCallback = fn
	}
}

// FrameBeginCallback 帧更新开始时的回调函数
func (WithOption) FrameBeginCallback(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameBeginCallback = fn
	}
}

// FrameEndCallback 帧更新结束时的回调函数
func (WithOption) FrameEndCallback(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameEndCallback = fn
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能
func (WithOption) FaceAnyAllocator(allocator container.Allocator[util.FaceAny]) ContextOption {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.FaceAnyAllocator = allocator
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能
func (WithOption) HookAllocator(allocator container.Allocator[localevent.Hook]) ContextOption {
	return func(o *ContextOptions) {
		if allocator == nil {
			panic("nil allocator")
		}
		o.HookAllocator = allocator
	}
}
