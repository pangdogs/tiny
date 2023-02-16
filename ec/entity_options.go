package ec

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// EntityOptions 创建实体的所有选项
type EntityOptions struct {
	Inheritor util.Face[Entity]                // 继承者，在扩展实体自身能力时使用
	FaceCache container.Cache[util.FaceAny]    // FaceCache用于提高性能，通常传入运行时上下文选项中的FaceCache
	HookCache container.Cache[localevent.Hook] // HookCache用于提高性能，通常传入运行时上下文选项中的HookCache
}

// EntityOption 创建实体的选项设置器
type EntityOption func(o *EntityOptions)

// WithEntityOption 创建实体的所有选项设置器
type WithEntityOption struct{}

// Default 默认值
func (WithEntityOption) Default() EntityOption {
	return func(o *EntityOptions) {
		WithEntityOption{}.Inheritor(util.Face[Entity]{})(o)
		WithEntityOption{}.FaceCache(nil)(o)
		WithEntityOption{}.HookCache(nil)(o)
	}
}

// Inheritor 继承者，在扩展实体自身能力时使用
func (WithEntityOption) Inheritor(v util.Face[Entity]) EntityOption {
	return func(o *EntityOptions) {
		o.Inheritor = v
	}
}

// FaceCache FaceCache用于提高性能，通常传入运行时上下文选项中的FaceCache
func (WithEntityOption) FaceCache(v container.Cache[util.FaceAny]) EntityOption {
	return func(o *EntityOptions) {
		o.FaceCache = v
	}
}

// HookCache HookCache用于提高性能，通常传入运行时上下文选项中的HookCache
func (WithEntityOption) HookCache(v container.Cache[localevent.Hook]) EntityOption {
	return func(o *EntityOptions) {
		o.HookCache = v
	}
}
