package runtime

import (
	"git.golaxy.org/tiny/utils/option"
	"time"
)

// NewFrame 创建帧，在运行时初始化时可以设置帧，用于设置运行时帧更新方式，在逻辑运行过程中可以在运行时上下文中获取帧信息。
func NewFrame(settings ...option.Setting[FrameOptions]) Frame {
	frame := &_FrameBehavior{}
	frame.init(option.Make(With.Frame.Default(), settings...))
	return frame
}

// Frame 帧，在运行时初始化时可以设置帧，用于设置运行时帧更新方式，在逻辑运行过程中可以在运行时上下文中获取帧信息
type Frame interface {
	iFrame

	// GetTargetFPS 获取目标FPS
	GetTargetFPS() float32
	// GetCurFPS 获取当前FPS
	GetCurFPS() float32
	// GetMode 获取帧模式
	GetMode() FrameMode
	// GetTotalFrames 获取运行帧数上限
	GetTotalFrames() int64
	// GetCurFrames 获取当前帧数
	GetCurFrames() int64
	// GetRunningBeginTime 获取运行开始时间
	GetRunningBeginTime() time.Time
	// GetRunningElapseTime 获取运行持续时间
	GetRunningElapseTime() time.Duration
	// GetLoopBeginTime 获取当前帧循环开始时间（包含异步调用）
	GetLoopBeginTime() time.Time
	// GetLastLoopElapseTime 获取上一帧循环耗时（包含异步调用）
	GetLastLoopElapseTime() time.Duration
	// GetUpdateBeginTime 获取当前帧更新开始时间
	GetUpdateBeginTime() time.Time
	// GetLastUpdateElapseTime 获取上一次帧更新耗时
	GetLastUpdateElapseTime() time.Duration
}

type iFrame interface {
	setCurFrames(v int64)
	runningBegin()
	runningEnd()
	loopBegin()
	loopEnd()
	updateBegin()
	updateEnd()
}

type _FrameBehavior struct {
	options              FrameOptions
	curFPS               float32
	curFrames            int64
	runningBeginTime     time.Time
	runningElapseTime    time.Duration
	loopBeginTime        time.Time
	lastLoopElapseTime   time.Duration
	updateBeginTime      time.Time
	lastUpdateElapseTime time.Duration
	statFPSBeginTime     time.Time
	statFPSFrames        int64
	fixedLoopElapseTime  time.Duration
}

// GetTargetFPS 获取目标FPS
func (frame *_FrameBehavior) GetTargetFPS() float32 {
	return frame.options.TargetFPS
}

// GetCurFPS 获取当前FPS
func (frame *_FrameBehavior) GetCurFPS() float32 {
	return frame.curFPS
}

// GetMode 获取帧模式
func (frame *_FrameBehavior) GetMode() FrameMode {
	return frame.options.Mode
}

// GetTotalFrames 获取运行帧数上限
func (frame *_FrameBehavior) GetTotalFrames() int64 {
	return frame.options.TotalFrames
}

// GetCurFrames 获取当前帧数
func (frame *_FrameBehavior) GetCurFrames() int64 {
	return frame.curFrames
}

// GetRunningBeginTime 获取运行开始时间
func (frame *_FrameBehavior) GetRunningBeginTime() time.Time {
	return frame.runningBeginTime
}

// GetRunningElapseTime 获取运行持续时间
func (frame *_FrameBehavior) GetRunningElapseTime() time.Duration {
	return frame.runningElapseTime
}

// GetLoopBeginTime 获取当前帧循环开始时间（包含异步调用）
func (frame *_FrameBehavior) GetLoopBeginTime() time.Time {
	return frame.loopBeginTime
}

// GetLastLoopElapseTime 获取上一帧循环耗时（包含异步调用）
func (frame *_FrameBehavior) GetLastLoopElapseTime() time.Duration {
	return frame.lastLoopElapseTime
}

// GetUpdateBeginTime 获取当前帧更新开始时间
func (frame *_FrameBehavior) GetUpdateBeginTime() time.Time {
	return frame.updateBeginTime
}

// GetLastUpdateElapseTime 获取上一次帧更新耗时
func (frame *_FrameBehavior) GetLastUpdateElapseTime() time.Duration {
	return frame.lastUpdateElapseTime
}

func (frame *_FrameBehavior) init(opts FrameOptions) {
	frame.options = opts
}

func (frame *_FrameBehavior) setCurFrames(v int64) {
	frame.curFrames = v
}

func (frame *_FrameBehavior) runningBegin() {
	now := time.Now()

	frame.curFPS = frame.options.TargetFPS
	frame.curFrames = 0

	frame.statFPSBeginTime = now
	frame.statFPSFrames = 0

	frame.runningBeginTime = now
	frame.runningElapseTime = 0

	frame.loopBeginTime = now
	frame.lastLoopElapseTime = 0

	frame.updateBeginTime = now
	frame.lastUpdateElapseTime = 0

	switch frame.options.Mode {
	case Simulate, Manual:
		frame.fixedLoopElapseTime = time.Second / time.Duration(frame.options.TargetFPS)
	}
}

func (frame *_FrameBehavior) runningEnd() {}

func (frame *_FrameBehavior) loopBegin() {
	switch frame.options.Mode {
	case Simulate, Manual:
		frame.loopBeginTime = frame.runningBeginTime.Add(frame.fixedLoopElapseTime * time.Duration(frame.curFrames))

	default:
		frame.loopBeginTime = time.Now()

		statInterval := frame.loopBeginTime.Sub(frame.statFPSBeginTime).Seconds()
		if statInterval >= 1 {
			frame.curFPS = float32(float64(frame.statFPSFrames) / statInterval)
			frame.statFPSBeginTime = frame.loopBeginTime
			frame.statFPSFrames = 0
		}
	}
}

func (frame *_FrameBehavior) loopEnd() {
	switch frame.options.Mode {
	case Simulate, Manual:
		frame.lastLoopElapseTime = frame.fixedLoopElapseTime
		frame.runningElapseTime += frame.fixedLoopElapseTime

	default:
		frame.lastLoopElapseTime = time.Now().Sub(frame.loopBeginTime)
		frame.runningElapseTime += frame.lastLoopElapseTime
		frame.statFPSFrames++
	}
}

func (frame *_FrameBehavior) updateBegin() {
	switch frame.options.Mode {
	case Simulate, Manual:
		frame.updateBeginTime = frame.runningBeginTime.Add(frame.fixedLoopElapseTime * time.Duration(frame.curFrames))
	default:
		frame.updateBeginTime = time.Now()
	}
}

func (frame *_FrameBehavior) updateEnd() {
	switch frame.options.Mode {
	case Simulate, Manual:
		frame.lastUpdateElapseTime = frame.fixedLoopElapseTime
	default:
		frame.lastUpdateElapseTime = time.Now().Sub(frame.updateBeginTime)
	}
}
