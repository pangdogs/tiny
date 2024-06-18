package tiny

import (
	"context"
	"fmt"
	"git.golaxy.org/tiny/event"
	"git.golaxy.org/tiny/internal/gctx"
	"git.golaxy.org/tiny/plugin"
	"git.golaxy.org/tiny/runtime"
	"git.golaxy.org/tiny/utils/generic"
	"time"
)

var (
	ErrCtrlChanClosed = fmt.Errorf("%w: ctrl chan is closed", ErrRuntime) // 运行控制已关闭
)

type _Ctrl struct {
	at     bool
	frames int64
}

// Run 运行
func (rt *RuntimeBehavior) Run() <-chan struct{} {
	ctx := rt.ctx

	select {
	case <-ctx.Done():
		panic(fmt.Errorf("%w: %w", ErrRuntime, context.Canceled))
	case <-ctx.TerminatedChan():
		panic(fmt.Errorf("%w: terminated", ErrRuntime))
	default:
	}

	if parentCtx, ok := ctx.GetParentContext().(gctx.Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	go rt.running()

	return gctx.UnsafeContext(ctx).GetTerminatedChan()
}

// Play 播放指定时长
func (rt *RuntimeBehavior) Play(delta time.Duration) (err error) {
	frame := rt.opts.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	frames := int64(delta.Seconds() * float64(frame.GetTargetFPS()))
	if frames <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		at:     false,
		frames: frames,
	}

	return nil
}

// PlayAt 播放至指定位置
func (rt *RuntimeBehavior) PlayAt(at time.Duration) (err error) {
	frame := rt.opts.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	frames := int64(at.Seconds() * float64(frame.GetTargetFPS()))
	if frames <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		at:     false,
		frames: frames,
	}

	return nil
}

// PlayFrames 播放指定帧数
func (rt *RuntimeBehavior) PlayFrames(delta int64) (err error) {
	frame := rt.opts.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	if delta <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		at:     false,
		frames: delta,
	}

	return nil
}

// PlayAtFrames 播放至指定帧数
func (rt *RuntimeBehavior) PlayAtFrames(at int64) (err error) {
	frame := rt.opts.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	if at <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		at:     true,
		frames: at,
	}

	return nil
}

// Terminate 停止
func (rt *RuntimeBehavior) Terminate() <-chan struct{} {
	return rt.ctx.Terminate()
}

// TerminatedChan 已停止chan
func (rt *RuntimeBehavior) TerminatedChan() <-chan struct{} {
	return rt.ctx.TerminatedChan()
}

func (rt *RuntimeBehavior) running() {
	ctx := rt.ctx

	rt.changeRunningState(runtime.RunningState_Starting)

	hooks := rt.loopStart()

	rt.changeRunningState(runtime.RunningState_Started)

	rt.mainLoop()

	rt.changeRunningState(runtime.RunningState_Terminating)

	rt.loopStop(hooks)
	ctx.GetWaitGroup().Wait()

	rt.changeRunningState(runtime.RunningState_Terminated)

	if parentCtx, ok := ctx.GetParentContext().(gctx.Context); ok {
		parentCtx.GetWaitGroup().Done()
	}

	close(gctx.UnsafeContext(ctx).GetTerminatedChan())
}

func (rt *RuntimeBehavior) changeRunningState(state runtime.RunningState) {
	switch state {
	case runtime.RunningState_Starting:
		rt.initPlugin()
	case runtime.RunningState_FrameLoopBegin:
		runtime.UnsafeFrame(rt.opts.Frame).LoopBegin()
	case runtime.RunningState_FrameUpdateBegin:
		runtime.UnsafeFrame(rt.opts.Frame).UpdateBegin()
	case runtime.RunningState_FrameUpdateEnd:
		runtime.UnsafeFrame(rt.opts.Frame).UpdateEnd()
	case runtime.RunningState_FrameLoopEnd:
		runtime.UnsafeFrame(rt.opts.Frame).LoopEnd()
	case runtime.RunningState_Terminated:
		rt.shutPlugin()
	}

	runtime.UnsafeContext(rt.ctx).ChangeRunningState(state)
}

func (rt *RuntimeBehavior) initPlugin() {
	pluginBundle := rt.ctx.GetPluginBundle()
	if pluginBundle == nil {
		return
	}

	plugin.UnsafePluginBundle(pluginBundle).SetInstallCB(rt.activatePlugin)
	plugin.UnsafePluginBundle(pluginBundle).SetUninstallCB(rt.deactivatePlugin)

	pluginBundle.Range(func(pluginInfo plugin.PluginInfo) bool {
		rt.activatePlugin(pluginInfo)
		return true
	})
}

func (rt *RuntimeBehavior) shutPlugin() {
	pluginBundle := rt.ctx.GetPluginBundle()
	if pluginBundle == nil {
		return
	}

	plugin.UnsafePluginBundle(pluginBundle).SetInstallCB(nil)
	plugin.UnsafePluginBundle(pluginBundle).SetUninstallCB(nil)

	pluginBundle.ReversedRange(func(pluginInfo plugin.PluginInfo) bool {
		rt.deactivatePlugin(pluginInfo)
		return true
	})
}

func (rt *RuntimeBehavior) activatePlugin(pluginInfo plugin.PluginInfo) {
	if pluginInit, ok := pluginInfo.Face.Iface.(LifecyclePluginInit); ok {
		generic.MakeAction1(pluginInit.InitP).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError(), rt.ctx)
	}
	plugin.UnsafePluginBundle(rt.ctx.GetPluginBundle()).SetActive(pluginInfo.Name, true)
}

func (rt *RuntimeBehavior) deactivatePlugin(pluginInfo plugin.PluginInfo) {
	plugin.UnsafePluginBundle(rt.ctx.GetPluginBundle()).SetActive(pluginInfo.Name, false)
	if pluginShut, ok := pluginInfo.Face.Iface.(LifecyclePluginShut); ok {
		generic.MakeAction1(pluginShut.ShutP).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError(), rt.ctx)
	}
}

func (rt *RuntimeBehavior) loopStart() (hooks [5]event.Hook) {
	ctx := rt.ctx
	frame := rt.opts.Frame

	if frame != nil {
		runtime.UnsafeFrame(frame).RunningBegin()
	}

	hooks[0] = runtime.BindEventEntityMgrAddEntity(ctx.GetEntityMgr(), rt)
	hooks[1] = runtime.BindEventEntityMgrRemoveEntity(ctx.GetEntityMgr(), rt)
	hooks[2] = runtime.BindEventEntityMgrEntityAddComponents(ctx.GetEntityMgr(), rt)
	hooks[3] = runtime.BindEventEntityMgrEntityRemoveComponent(ctx.GetEntityMgr(), rt)
	hooks[4] = runtime.BindEventEntityMgrEntityFirstAccessComponent(ctx.GetEntityMgr(), rt)

	return
}

func (rt *RuntimeBehavior) loopStop(hooks [5]event.Hook) {
	frame := rt.opts.Frame

	event.Clean(hooks[:])

	if frame != nil {
		runtime.UnsafeFrame(frame).RunningEnd()
	}
}

func (rt *RuntimeBehavior) mainLoop() {
	frame := rt.opts.Frame

	if frame == nil {
		rt.loopingNoFrame()
	} else {
		switch frame.GetMode() {
		case runtime.Simulate:
			rt.loopingSimulate()
		case runtime.Manual:
			rt.loopingManual()
		default:
			rt.loopingRealTime()
		}
	}
}

func (rt *RuntimeBehavior) runTask(task _Task) {
	switch task.typ {
	case _TaskType_Call:
		rt.changeRunningState(runtime.RunningState_RunCallBegin)
		task.run(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		rt.changeRunningState(runtime.RunningState_RunCallEnd)
	case _TaskType_Frame:
		task.run(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}
}

func (rt *RuntimeBehavior) runGC() {
	rt.changeRunningState(runtime.RunningState_RunGCBegin)
	rt.gc()
	rt.changeRunningState(runtime.RunningState_RunGCEnd)
}
