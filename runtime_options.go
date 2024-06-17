package tiny

import (
	"fmt"
	"git.golaxy.org/tiny/runtime"
	"git.golaxy.org/tiny/utils/generic"
	"git.golaxy.org/tiny/utils/iface"
	"git.golaxy.org/tiny/utils/option"
	"time"
)

type (
	CustomGC = generic.DelegateAction1[Runtime] // 自定义GC函数
)

// RuntimeOptions 创建运行时的所有选项
type RuntimeOptions struct {
	CompositeFace        iface.Face[Runtime] // 扩展者，在扩展运行时自身能力时使用
	AutoRun              bool                // 是否开启自动运行
	ProcessQueueCapacity int                 // 任务处理流水线大小
	Frame                runtime.Frame       // 帧，设置为nil表示不使用帧更新特性
	GCInterval           time.Duration       // GC间隔时长
	CustomGC             CustomGC            // 自定义GC
}

var With _RuntimeOption

type _RuntimeOption struct{}

// Default 运行时的默认值
func (_RuntimeOption) Default() option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		With.CompositeFace(iface.Face[Runtime]{})(o)
		With.AutoRun(false)(o)
		With.ProcessQueueCapacity(128)(o)
		With.Frame(nil)(o)
		With.GCInterval(10 * time.Second)(o)
		With.CustomGC(nil)(o)
	}
}

// CompositeFace 扩展者，在扩展运行时自身能力时使用
func (_RuntimeOption) CompositeFace(face iface.Face[Runtime]) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		o.CompositeFace = face
	}
}

// AutoRun 运行时是否开启自动运行
func (_RuntimeOption) AutoRun(b bool) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		o.AutoRun = b
	}
}

// ProcessQueueCapacity 任务处理流水线大小
func (_RuntimeOption) ProcessQueueCapacity(cap int) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		if cap <= 0 {
			panic(fmt.Errorf("%w: %w: ProcessQueueCapacity less equal 0 is invalid", ErrRuntime, ErrArgs))
		}
		o.ProcessQueueCapacity = cap
	}
}

// Frame 运行时的帧，设置为nil表示不使用帧更新特性
func (_RuntimeOption) Frame(frame runtime.Frame) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		o.Frame = frame
	}
}

// GCInterval 运行时的GC间隔时长
func (_RuntimeOption) GCInterval(dur time.Duration) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		if dur <= 0 {
			panic(fmt.Errorf("%w: %w: GCInterval less equal 0 is invalid", ErrRuntime, ErrArgs))
		}
		o.GCInterval = dur
	}
}

// CustomGC 运行时的自定义GC
func (_RuntimeOption) CustomGC(fn CustomGC) option.Setting[RuntimeOptions] {
	return func(o *RuntimeOptions) {
		o.CustomGC = fn
	}
}
