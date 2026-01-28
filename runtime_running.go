package tiny

import (
	"context"
	"fmt"
	"time"

	"git.golaxy.org/core/event"
	"git.golaxy.org/core/extension"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/corectx"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/tiny/ec"
	"git.golaxy.org/tiny/ec/pt"
	"git.golaxy.org/tiny/runtime"
	"git.golaxy.org/tiny/utils/exception"
)

var (
	ErrCtrlChanClosed = fmt.Errorf("%w: ctrl chan is closed", ErrRuntime) // 运行控制已关闭
)

type _CtrlMode int32

const (
	_NoCtrl _CtrlMode = iota
	_FrameDelta
	_FrameAt
	_FuncAt
)

type _Ctrl struct {
	mode   _CtrlMode
	frames int64
	fun    generic.Func1[runtime.Context, bool]
}

// Run 运行
func (rt *RuntimeBehavior) Run() async.AsyncRet {
	ctx := rt.ctx

	select {
	case <-ctx.Done():
		exception.Panicf("%w: %w", ErrRuntime, context.Canceled)
	case <-ctx.Terminated():
		exception.Panicf("%w: terminated", ErrRuntime)
	default:
	}

	if !rt.isRunning.CompareAndSwap(false, true) {
		exception.Panicf("%w: already running", ErrRuntime)
	}

	if parentCtx, ok := ctx.GetParentContext().(corectx.Context); ok {
		parentCtx.GetWaitGroup().Add(1)
	}

	go rt.running()

	return ctx.Terminated()
}

// Play 播放指定时长
func (rt *RuntimeBehavior) Play(delta time.Duration) (err error) {
	frame := rt.options.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	frames := int64(delta.Seconds() * frame.GetTargetFPS())
	if frames <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		mode:   _FrameDelta,
		frames: frames,
	}

	return nil
}

// PlayAt 播放至指定位置
func (rt *RuntimeBehavior) PlayAt(at time.Duration) (err error) {
	frame := rt.options.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	frames := int64(at.Seconds() * frame.GetTargetFPS())
	if frames <= 0 {
		return nil
	}

	rt.ctrlChan <- _Ctrl{
		mode:   _FrameAt,
		frames: frames,
	}

	return nil
}

// PlayFrames 播放指定帧数
func (rt *RuntimeBehavior) PlayFrames(delta int64) (err error) {
	frame := rt.options.Frame

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
		mode:   _FrameDelta,
		frames: delta,
	}

	return nil
}

// PlayAtFrames 播放至指定帧数
func (rt *RuntimeBehavior) PlayAtFrames(at int64) (err error) {
	frame := rt.options.Frame

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
		mode:   _FrameAt,
		frames: at,
	}

	return nil
}

// PlayAtFunc 播放至函数指定位置
func (rt *RuntimeBehavior) PlayAtFunc(fun generic.Func1[runtime.Context, bool]) (err error) {
	frame := rt.options.Frame

	if frame == nil || frame.GetMode() != runtime.Manual {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	rt.ctrlChan <- _Ctrl{
		mode: _FuncAt,
		fun:  fun,
	}

	return nil
}

// Terminate 停止
func (rt *RuntimeBehavior) Terminate() async.AsyncRet {
	return rt.ctx.Terminate()
}

// Terminated 已停止
func (rt *RuntimeBehavior) Terminated() async.AsyncRet {
	return rt.ctx.Terminated()
}

func (rt *RuntimeBehavior) running() {
	ctx := rt.ctx

	rt.emitEventRunningEvent(runtime.RunningEvent_Starting)

	handles := rt.loopStart()

	rt.emitEventRunningEvent(runtime.RunningEvent_Started)

	rt.mainLoop()

	rt.emitEventRunningEvent(runtime.RunningEvent_Terminating)

	rt.loopStop(handles)
	ctx.GetWaitGroup().Wait()

	rt.emitEventRunningEvent(runtime.RunningEvent_Terminated)

	if parentCtx, ok := ctx.GetParentContext().(corectx.Context); ok {
		parentCtx.GetWaitGroup().Done()
	}

	corectx.UnsafeContext(ctx).ReturnTerminated()
}

func (rt *RuntimeBehavior) emitEventRunningEvent(runningEvent runtime.RunningEvent, args ...any) {
	runtime.UnsafeContext(rt.ctx).EmitEventRunningEvent(runningEvent, args...)
}

func (rt *RuntimeBehavior) onBeforeContextRunningEvent(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
	switch runningEvent {
	case runtime.RunningEvent_Birth:
		if rt.options.AutoRun {
			rt.getInstance().Run()
		}
	case runtime.RunningEvent_Starting:
		rt.initComponentPT()
		rt.initEntityPT()
		rt.initAddIn()
	case runtime.RunningEvent_FrameLoopBegin:
		runtime.UnsafeFrame(rt.options.Frame).LoopBegin()
	case runtime.RunningEvent_FrameUpdateBegin:
		runtime.UnsafeFrame(rt.options.Frame).UpdateBegin()
	case runtime.RunningEvent_FrameUpdateEnd:
		runtime.UnsafeFrame(rt.options.Frame).UpdateEnd()
	case runtime.RunningEvent_FrameLoopEnd:
		runtime.UnsafeFrame(rt.options.Frame).LoopEnd()
	}
}

func (rt *RuntimeBehavior) onAfterContextRunningEvent(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
	switch runningEvent {
	case runtime.RunningEvent_Terminated:
		rt.shutAddIn()
		rt.shutEntityPT()
		rt.shutComponentPT()
	}
}

func (rt *RuntimeBehavior) initEntityPT() {
	rt.managedEntityLibHandles[0] = pt.BindEventEntityLibDeclareEntityPT(rt.ctx.GetEntityLib(), pt.HandleEventEntityLibDeclareEntityPT(rt.onEntityLibDeclareEntityPT))

	for _, entityPT := range rt.ctx.GetEntityLib().List() {
		rt.emitEventRunningEvent(runtime.RunningEvent_EntityPTDeclared, entityPT)
	}
}

func (rt *RuntimeBehavior) shutEntityPT() {
	event.UnbindHandles(rt.managedEntityLibHandles[:])
}

func (rt *RuntimeBehavior) initComponentPT() {
	rt.managedComponentLibHandles[0] = pt.BindEventComponentLibDeclareComponentPT(rt.ctx.GetEntityLib().GetComponentLib(), pt.HandleEventComponentLibDeclareComponentPT(rt.onComponentLibDeclareComponentPT))

	for _, compPT := range rt.ctx.GetEntityLib().GetComponentLib().List() {
		rt.emitEventRunningEvent(runtime.RunningEvent_ComponentPTDeclared, compPT)
	}
}

func (rt *RuntimeBehavior) shutComponentPT() {
	event.UnbindHandles(rt.managedComponentLibHandles[:])
}

func (rt *RuntimeBehavior) onEntityLibDeclareEntityPT(entityPT ec.EntityPT) {
	rt.emitEventRunningEvent(runtime.RunningEvent_EntityPTDeclared, entityPT)
}

func (rt *RuntimeBehavior) onComponentLibDeclareComponentPT(compPT ec.ComponentPT) {
	rt.emitEventRunningEvent(runtime.RunningEvent_ComponentPTDeclared, compPT)
}

func (rt *RuntimeBehavior) initAddIn() {
	addInManager := runtime.UnsafeContext(rt.ctx).GetAddInManager()

	rt.managedAddInManagerHandles[0] = extension.BindEventRuntimeInstallAddIn(addInManager, extension.HandleEventRuntimeInstallAddIn(rt.activateAddIn))
	rt.managedAddInManagerHandles[1] = extension.BindEventRuntimeUninstallAddIn(addInManager, extension.HandleEventRuntimeUninstallAddIn(rt.deactivateAddIn))

	addInStatusList := addInManager.List()
	for i := range addInStatusList {
		rt.activateAddIn(addInStatusList[i])
	}
}

func (rt *RuntimeBehavior) shutAddIn() {
	addInManager := runtime.UnsafeContext(rt.ctx).GetAddInManager()

	rt.managedAddInManagerHandles[0].Unbind()

	addInStatusList := addInManager.List()
	for i := len(addInStatusList) - 1; i >= 0; i-- {
		addInStatusList[i].Uninstall()
	}

	rt.managedAddInManagerHandles[1].Unbind()
}

func (rt *RuntimeBehavior) activateAddIn(status extension.AddInStatus) {
	if status.State() != extension.AddInState_Loaded {
		return
	}

	rt.emitEventRunningEvent(runtime.RunningEvent_AddInActivating, status)

	if status.State() != extension.AddInState_Loaded {
		rt.emitEventRunningEvent(runtime.RunningEvent_AddInActivatingAborted, status)
		return
	}

	if cb, ok := status.InstanceFace().Iface.(LifecycleRuntimeAddInInit); ok {
		generic.CastAction1(cb.Init).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError(), rt.ctx)
	}

	if status.State() != extension.AddInState_Loaded {
		rt.emitEventRunningEvent(runtime.RunningEvent_AddInActivatingAborted, status)
		return
	}

	addInStatus := status.(extension.RuntimeAddInStatus)
	extension.UnsafeRuntimeAddInStatus(addInStatus).SetState(extension.AddInState_Running)

	rt.emitEventRunningEvent(runtime.RunningEvent_AddInActivatingDone, status)

	if status.State() != extension.AddInState_Running {
		return
	}

	if cb, ok := status.InstanceFace().Iface.(LifecycleAddInOnRuntimeRunningEvent); ok {
		extension.UnsafeRuntimeAddInStatus(addInStatus).ManagedRuntimeRunningEventHandle(
			runtime.BindEventContextRunningEvent(rt.ctx, runtime.HandleEventContextRunningEvent(cb.OnContextRunningEvent)),
		)
	}
}

func (rt *RuntimeBehavior) deactivateAddIn(status extension.AddInStatus) {
	if status.State() != extension.AddInState_Running {
		return
	}

	rt.emitEventRunningEvent(runtime.RunningEvent_AddInDeactivating, status)

	if cb, ok := status.InstanceFace().Iface.(LifecycleRuntimeAddInShut); ok {
		generic.CastAction1(cb.Shut).Call(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError(), rt.ctx)
	}

	rt.emitEventRunningEvent(runtime.RunningEvent_AddInDeactivatingDone, status)
}

func (rt *RuntimeBehavior) loopStart() []event.Handle {
	ctx := rt.ctx
	frame := rt.options.Frame

	if frame != nil {
		runtime.UnsafeFrame(frame).RunningBegin()
	}

	return []event.Handle{
		runtime.BindEventEntityManagerAddEntity(ctx.GetEntityManager(), rt.handleEventEntityManagerAddEntity),
		runtime.BindEventEntityManagerRemoveEntity(ctx.GetEntityManager(), rt.handleEventEntityManagerRemoveEntity),
		runtime.BindEventEntityManagerEntityAddComponents(ctx.GetEntityManager(), rt.handleEventEntityManagerEntityAddComponents),
		runtime.BindEventEntityManagerEntityRemoveComponent(ctx.GetEntityManager(), rt.handleEventEntityManagerEntityRemoveComponent),
		runtime.BindEventEntityManagerEntityComponentEnableChanged(ctx.GetEntityManager(), rt.handleEventEntityManagerEntityComponentEnableChanged),
		runtime.BindEventEntityManagerEntityFirstTouchComponent(ctx.GetEntityManager(), rt.handleEventEntityManagerEntityFirstTouchComponent),
	}
}

func (rt *RuntimeBehavior) loopStop(handles []event.Handle) {
	frame := rt.options.Frame

	event.UnbindHandles(handles)

	if frame != nil {
		runtime.UnsafeFrame(frame).RunningEnd()
	}
}

func (rt *RuntimeBehavior) mainLoop() {
	frame := rt.options.Frame

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
		rt.emitEventRunningEvent(runtime.RunningEvent_RunCallBegin)
		task.run(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
		rt.emitEventRunningEvent(runtime.RunningEvent_RunCallEnd)
	case _TaskType_Frame:
		task.run(rt.ctx.GetAutoRecover(), rt.ctx.GetReportError())
	}
}

func (rt *RuntimeBehavior) runGC() {
	rt.emitEventRunningEvent(runtime.RunningEvent_RunGCBegin)
	rt.gc()
	rt.emitEventRunningEvent(runtime.RunningEvent_RunGCEnd)
}
