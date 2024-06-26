package tiny

import (
	"git.golaxy.org/tiny/runtime"
	"git.golaxy.org/tiny/utils/generic"
	"time"
)

// Running 运行接口
type Running interface {
	// Run 运行
	Run() <-chan struct{}
	// Play 播放指定时长
	Play(delta time.Duration) error
	// PlayAt 播放至指定位置
	PlayAt(at time.Duration) error
	// PlayFrames 播放指定帧数
	PlayFrames(delta int64) error
	// PlayAtFrames 播放至指定帧数
	PlayAtFrames(at int64) error
	// PlayAtFunc 播放至函数指定位置
	PlayAtFunc(fun generic.Func1[runtime.Context, bool]) error
	// Terminate 停止
	Terminate() <-chan struct{}
	// TerminatedChan 已停止chan
	TerminatedChan() <-chan struct{}
}
