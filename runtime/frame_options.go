package runtime

import (
	"fmt"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/option"
)

// FrameMode 帧模式
type FrameMode int32

const (
	RealTime FrameMode = iota // 实时
	Simulate                  // 瞬时模拟
	Manual                    // 手动控制
)

// FrameOptions 帧的所有选项
type FrameOptions struct {
	TargetFPS   float32   // 目标FPS
	TotalFrames int64     // 运行帧数上限
	Mode        FrameMode // 帧模式
}

type _FrameOption struct{}

// Default 默认值
func (_FrameOption) Default() option.Setting[FrameOptions] {
	return func(o *FrameOptions) {
		With.Frame.TargetFPS(30)(o)
		With.Frame.TotalFrames(0)(o)
		With.Frame.Mode(RealTime)(o)
	}
}

// TargetFPS 目标FPS
func (_FrameOption) TargetFPS(fps float32) option.Setting[FrameOptions] {
	return func(o *FrameOptions) {
		if fps <= 0 {
			panic(fmt.Errorf("%w: %w: TargetFPS less equal 0 is invalid", ErrFrame, exception.ErrArgs))
		}
		o.TargetFPS = fps
	}
}

// TotalFrames 运行帧数上限
func (_FrameOption) TotalFrames(v int64) option.Setting[FrameOptions] {
	return func(o *FrameOptions) {
		if v < 0 {
			panic(fmt.Errorf("%w: %w: TotalFrames less 0 is invalid", ErrFrame, exception.ErrArgs))
		}
		o.TotalFrames = v
	}
}

// Mode 帧模式
func (_FrameOption) Mode(m FrameMode) option.Setting[FrameOptions] {
	return func(o *FrameOptions) {
		o.Mode = m
	}
}
