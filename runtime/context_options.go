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
	StartedCb          func(ctx Context)                    // 启动运行时回调函数
	StoppingCb         func(ctx Context)                    // 开始停止运行时回调函数
	StoppedCb          func(ctx Context)                    // 完全停止运行时回调函数
	FrameBeginCb       func(ctx Context)                    // 帧开始时的回调函数
	FrameEndCb         func(ctx Context)                    // 帧结束时的回调函数
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
		WithOption{}.StartedCb(nil)(o)
		WithOption{}.StoppingCb(nil)(o)
		WithOption{}.StoppedCb(nil)(o)
		WithOption{}.FrameBeginCb(nil)(o)
		WithOption{}.FrameEndCb(nil)(o)
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

// StartedCb 启动运行时回调函数
func (WithOption) StartedCb(v func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StartedCb = v
	}
}

// StoppingCb 开始停止运行时回调函数
func (WithOption) StoppingCb(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppingCb = fn
	}
}

// StoppedCb 完全停止运行时回调函数
func (WithOption) StoppedCb(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.StoppedCb = fn
	}
}

// FrameBeginCb 帧更新开始时的回调函数
func (WithOption) FrameBeginCb(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameBeginCb = fn
	}
}

// FrameEndCb 帧更新结束时的回调函数
func (WithOption) FrameEndCb(fn func(ctx Context)) ContextOption {
	return func(o *ContextOptions) {
		o.FrameEndCb = fn
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
