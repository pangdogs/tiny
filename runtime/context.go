package runtime

import (
	"context"
	"fmt"
	"git.golaxy.org/tiny/event"
	"git.golaxy.org/tiny/internal/gctx"
	"git.golaxy.org/tiny/plugin"
	"git.golaxy.org/tiny/utils/async"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/iface"
	"git.golaxy.org/tiny/utils/option"
	"git.golaxy.org/tiny/utils/pool"
	"git.golaxy.org/tiny/utils/reinterpret"
	"git.golaxy.org/tiny/utils/uid"
	"reflect"
)

// NewContext 创建运行时上下文
func NewContext(settings ...option.Setting[ContextOptions]) Context {
	return UnsafeNewContext(option.Make(With.Context.Default(), settings...))
}

// Deprecated: UnsafeNewContext 内部创建运行时上下文
func UnsafeNewContext(options ContextOptions) Context {
	if !options.CompositeFace.IsNil() {
		options.CompositeFace.Iface.init(options)
		return options.CompositeFace.Iface
	}

	ctx := &ContextBehavior{}
	ctx.init(options)

	return ctx.opts.CompositeFace.Iface
}

// Context 运行时上下文接口
type Context interface {
	iContext
	gctx.CurrentContextProvider
	gctx.Context
	async.Caller
	reinterpret.CompositeProvider
	plugin.PluginProvider
	pool.ManagedPoolObject
	GCCollector

	// GetReflected 获取反射值
	GetReflected() reflect.Value
	// GetFrame 获取帧
	GetFrame() Frame
	// GetEntityMgr 获取实体管理器
	GetEntityMgr() EntityMgr
	// GetEntityTree 获取实体树
	GetEntityTree() EntityTree
	// ActivateEvent 启用事件
	ActivateEvent(event event.IEventCtrl, recursion event.EventRecursion)
	// ManagedHooks 托管hook，在运行时停止时自动解绑定
	ManagedHooks(hooks ...event.Hook)
	// AutoManagedPoolObject 自动判断托管对象池
	AutoManagedPoolObject() pool.ManagedPoolObject
}

type iContext interface {
	init(opts ContextOptions)
	getOptions() *ContextOptions
	newId() uid.Id
	setFrame(frame Frame)
	setCallee(callee async.Callee)
	changeRunningState(state RunningState)
	gc()
}

// ContextBehavior 运行时上下文行为，在需要扩展运行时上下文能力时，匿名嵌入至运行时上下文结构体中
type ContextBehavior struct {
	gctx.ContextBehavior
	opts               ContextOptions
	genId              uid.Id
	reflected          reflect.Value
	frame              Frame
	entityMgr          _EntityMgrBehavior
	callee             async.Callee
	managedHooks       []event.Hook
	managedPoolObjects []pool.PoolObject
	gcList             []GC
}

// GetReflected 获取反射值
func (ctx *ContextBehavior) GetReflected() reflect.Value {
	return ctx.reflected
}

// GetFrame 获取帧
func (ctx *ContextBehavior) GetFrame() Frame {
	return ctx.frame
}

// GetEntityMgr 获取实体管理器
func (ctx *ContextBehavior) GetEntityMgr() EntityMgr {
	return &ctx.entityMgr
}

// GetEntityTree 获取主实体树
func (ctx *ContextBehavior) GetEntityTree() EntityTree {
	return &ctx.entityMgr
}

// ActivateEvent 启用事件
func (ctx *ContextBehavior) ActivateEvent(event event.IEventCtrl, recursion event.EventRecursion) {
	if event == nil {
		panic(fmt.Errorf("%w: %w: event is nil", ErrContext, exception.ErrArgs))
	}
	event.Init(ctx.GetAutoRecover(), ctx.GetReportError(), recursion, ctx.AutoManagedPoolObject())
}

// GetCurrentContext 获取当前上下文
func (ctx *ContextBehavior) GetCurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.opts.CompositeFace.Iface)
}

// GetConcurrentContext 获取多线程安全的上下文
func (ctx *ContextBehavior) GetConcurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.opts.CompositeFace.Iface)
}

// GetCompositeFaceCache 支持重新解释类型
func (ctx *ContextBehavior) GetCompositeFaceCache() iface.Cache {
	return ctx.opts.CompositeFace.Cache
}

// CollectGC 收集GC
func (ctx *ContextBehavior) CollectGC(gc GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	ctx.gcList = append(ctx.gcList, gc)
}

func (ctx *ContextBehavior) init(opts ContextOptions) {
	ctx.opts = opts

	if ctx.opts.CompositeFace.IsNil() {
		ctx.opts.CompositeFace = iface.MakeFaceT[Context](ctx)
	}

	if ctx.opts.Context == nil {
		ctx.opts.Context = context.Background()
	}

	if ctx.opts.UseObjectPool {
		ctx.managedPoolObjects = make([]pool.PoolObject, 0, ctx.opts.UseObjectPoolSize)
	}

	gctx.UnsafeContext(&ctx.ContextBehavior).Init(ctx.opts.Context, ctx.opts.AutoRecover, ctx.opts.ReportError)
	ctx.reflected = reflect.ValueOf(ctx.opts.CompositeFace.Iface)
	ctx.entityMgr.init(ctx.opts.CompositeFace.Iface)
}

func (ctx *ContextBehavior) getOptions() *ContextOptions {
	return &ctx.opts
}

func (ctx *ContextBehavior) newId() uid.Id {
	ctx.genId++
	return ctx.genId
}

func (ctx *ContextBehavior) setFrame(frame Frame) {
	ctx.frame = frame
}

func (ctx *ContextBehavior) setCallee(callee async.Callee) {
	ctx.callee = callee
}

func (ctx *ContextBehavior) changeRunningState(state RunningState) {
	ctx.entityMgr.changeRunningState(state)
	ctx.opts.RunningHandler.Call(ctx.GetAutoRecover(), ctx.GetReportError(), nil, ctx.opts.CompositeFace.Iface, state)

	switch state {
	case RunningState_Terminated:
		ctx.cleanManagedHooks()
		ctx.cleanManagedPoolObjects()
	}
}
