package ec

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// EntityOptions 创建实体的所有选项
type EntityOptions struct {
	Inheritor        util.Face[Entity]                    // 继承者，在扩展实体自身能力时使用
	FaceAnyAllocator container.Allocator[util.FaceAny]    // 自定义FaceAny内存分配器，用于提高性能，通常传入运行时上下文中的FaceAnyAllocator
	HookAllocator    container.Allocator[localevent.Hook] // 自定义Hook内存分配器，用于提高性能，通常传入运行时上下文中的HookAllocator
	GCCollector      container.GCCollector                // 自定义GC收集器，通常不传或者传入运行时上下文
}

// EntityOption 创建实体的选项设置器
type EntityOption func(o *EntityOptions)

// WithEntityOption 创建实体的所有选项设置器
type WithEntityOption struct{}

// Default 默认值
func (WithEntityOption) Default() EntityOption {
	return func(o *EntityOptions) {
		WithEntityOption{}.Inheritor(util.Face[Entity]{})(o)
		WithEntityOption{}.FaceAnyAllocator(container.DefaultAllocator[util.FaceAny]())(o)
		WithEntityOption{}.HookAllocator(container.DefaultAllocator[localevent.Hook]())(o)
		WithEntityOption{}.GCCollector(nil)(o)
	}
}

// Inheritor 继承者，在扩展实体自身能力时使用
func (WithEntityOption) Inheritor(v util.Face[Entity]) EntityOption {
	return func(o *EntityOptions) {
		o.Inheritor = v
	}
}

// FaceAnyAllocator 自定义FaceAny内存分配器，用于提高性能，通常传入运行时上下文中的FaceAnyAllocator
func (WithEntityOption) FaceAnyAllocator(v container.Allocator[util.FaceAny]) EntityOption {
	return func(o *EntityOptions) {
		if v == nil {
			panic("nil allocator")
		}
		o.FaceAnyAllocator = v
	}
}

// HookAllocator 自定义Hook内存分配器，用于提高性能，通常传入运行时上下文中的HookAllocator
func (WithEntityOption) HookAllocator(v container.Allocator[localevent.Hook]) EntityOption {
	return func(o *EntityOptions) {
		if v == nil {
			panic("nil allocator")
		}
		o.HookAllocator = v
	}
}

// GCCollector 自定义GC收集器，通常不传或者传入运行时上下文
func (WithEntityOption) GCCollector(v container.GCCollector) EntityOption {
	return func(o *EntityOptions) {
		o.GCCollector = v
	}
}
