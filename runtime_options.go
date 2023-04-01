package tiny

import (
	"kit.golaxy.org/tiny/runtime"
	"kit.golaxy.org/tiny/util"
	"time"
)

// RuntimeOptions 创建运行时的所有选项
type RuntimeOptions struct {
	CompositeFace        util.Face[Runtime] // 扩展者，需要扩展运行时自身功能时需要使用
	EnableAutoRun        bool               // 是否开启自动运行
	ProcessQueueCapacity int                // 任务处理流水线大小
	ProcessQueueTimeout  time.Duration      // 当任务处理流水线满时，向其插入代码片段的超时时间，为0表示不等待直接报错
	SyncCallTimeout      time.Duration      // 同步调用超时时间，为0表示不处理超时，此时两个运行时互相同步调用会死锁
	Frame                runtime.Frame      // 帧
	GCInterval           time.Duration      // GC间隔时长
	CustomGC             func(rt Runtime)   // 自定义GC
}

// RuntimeOption 创建运行时的选项设置器
type RuntimeOption func(o *RuntimeOptions)

// WithRuntimeOption 创建运行时的所有选项设置器
type WithRuntimeOption struct{}

// Default 默认值
func (WithRuntimeOption) Default() RuntimeOption {
	return func(o *RuntimeOptions) {
		WithRuntimeOption{}.CompositeFace(util.Face[Runtime]{})(o)
		WithRuntimeOption{}.EnableAutoRun(false)(o)
		WithRuntimeOption{}.ProcessQueueCapacity(128)(o)
		WithRuntimeOption{}.ProcessQueueTimeout(0)(o)
		WithRuntimeOption{}.SyncCallTimeout(0)(o)
		WithRuntimeOption{}.Frame(nil)(o)
		WithRuntimeOption{}.GCInterval(10 * time.Second)(o)
		WithRuntimeOption{}.CustomGC(nil)(o)
	}
}

// CompositeFace 扩展者，需要扩展运行时自身功能时需要使用
func (WithRuntimeOption) CompositeFace(v util.Face[Runtime]) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.CompositeFace = v
	}
}

// EnableAutoRun 是否开启自动运行
func (WithRuntimeOption) EnableAutoRun(v bool) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.EnableAutoRun = v
	}
}

// ProcessQueueCapacity 任务处理流水线大小
func (WithRuntimeOption) ProcessQueueCapacity(v int) RuntimeOption {
	return func(o *RuntimeOptions) {
		if v <= 0 {
			panic("ProcessQueueCapacity less equal 0 is invalid")
		}
		o.ProcessQueueCapacity = v
	}
}

// ProcessQueueTimeout 当任务处理流水线满时，向其插入代码片段的超时时间，为0表示不等待直接报错
func (WithRuntimeOption) ProcessQueueTimeout(v time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.ProcessQueueTimeout = v
	}
}

// SyncCallTimeout 同步调用超时时间，为0表示不处理超时，此时两个运行时互相同步调用会死锁
func (WithRuntimeOption) SyncCallTimeout(v time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.SyncCallTimeout = v
	}
}

// Frame 帧
func (WithRuntimeOption) Frame(v runtime.Frame) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.Frame = v
	}
}

// GCInterval GC间隔时长
func (WithRuntimeOption) GCInterval(v time.Duration) RuntimeOption {
	return func(o *RuntimeOptions) {
		if v <= 0 {
			panic("GCInterval less equal 0 is invalid")
		}
		o.GCInterval = v
	}
}

// CustomGC 自定义GC
func (WithRuntimeOption) CustomGC(v func(rt Runtime)) RuntimeOption {
	return func(o *RuntimeOptions) {
		o.CustomGC = v
	}
}
