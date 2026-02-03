/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package runtime

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"

	"git.golaxy.org/core/event"
	"git.golaxy.org/core/extension"
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/corectx"
	"git.golaxy.org/core/utils/iface"
	"git.golaxy.org/core/utils/option"
	"git.golaxy.org/core/utils/reinterpret"
	"git.golaxy.org/tiny/ec/pt"
	"git.golaxy.org/tiny/utils/uid"
)

// NewContext 创建运行时上下文
func NewContext(settings ...option.Setting[ContextOptions]) Context {
	return UnsafeNewContext(option.New(With.Default(), settings...))
}

// Deprecated: UnsafeNewContext 内部创建运行时上下文
func UnsafeNewContext(options ContextOptions) Context {
	var ctx Context

	if !options.InstanceFace.IsNil() {
		ctx = options.InstanceFace.Iface
	} else {
		ctx = &ContextBehavior{}
	}
	ctx.init(options)

	return ctx
}

// Context 运行时上下文接口
type Context interface {
	iContext
	iConcurrentContext
	corectx.Context
	corectx.CurrentContextProvider
	reinterpret.InstanceProvider
	pt.EntityPTProvider
	extension.AddInProvider
	async.Caller
	GCCollector
	fmt.Stringer

	// GetName 获取名称
	GetName() string
	// GetReflected 获取反射值
	GetReflected() reflect.Value
	// GenUID 生成uid
	GenUID() uid.Id
	// GetFrame 获取帧
	GetFrame() Frame
	// GetEntityManager 获取实体管理器
	GetEntityManager() EntityManager
	// GetEntityTree 获取实体树
	GetEntityTree() EntityTree
	// Managed 托管事件句柄
	Managed() *event.ManagedHandles

	IContextRunningEventTab
}

type iContext interface {
	init(options ContextOptions)
	getOptions() *ContextOptions
	emitEventRunningEvent(runningEvent RunningEvent, args ...any)
	setFrame(frame Frame)
	setCallee(callee async.Callee)
	getAddInManager() extension.RuntimeAddInManager
	getScoped() *atomic.Bool
	gc()
}

// ContextBehavior 运行时上下文行为，在扩展运行时上下文能力时，匿名嵌入至运行时上下文结构体中
type ContextBehavior struct {
	corectx.ContextBehavior
	options       ContextOptions
	reflected     reflect.Value
	uidGenerator  int64
	frame         Frame
	entityManager _EntityManagerBehavior
	callee        async.Callee
	scoped        atomic.Bool
	gcList        []GC
	managed       event.ManagedHandles
	stringerOnce  sync.Once
	stringerCache string

	contextRunningEventTab contextRunningEventTab
}

// GetName 获取名称
func (ctx *ContextBehavior) GetName() string {
	return ctx.options.Name
}

// GetReflected 获取反射值
func (ctx *ContextBehavior) GetReflected() reflect.Value {
	return ctx.reflected
}

// GenUID 生成uid
func (ctx *ContextBehavior) GenUID() uid.Id {
	if ctx.options.UIDGenerator != nil {
		return uid.Id(ctx.options.UIDGenerator.Add(1))
	}
	ctx.uidGenerator++
	return uid.Id(ctx.uidGenerator)
}

// GetFrame 获取帧
func (ctx *ContextBehavior) GetFrame() Frame {
	return ctx.frame
}

// GetEntityManager 获取实体管理器
func (ctx *ContextBehavior) GetEntityManager() EntityManager {
	return &ctx.entityManager
}

// GetEntityTree 获取主实体树
func (ctx *ContextBehavior) GetEntityTree() EntityTree {
	return &ctx.entityManager
}

// Managed 托管事件句柄
func (ctx *ContextBehavior) Managed() *event.ManagedHandles {
	return &ctx.managed
}

// EventContextRunningEvent 事件：接收运行事件
func (ctx *ContextBehavior) EventContextRunningEvent() event.IEvent {
	return ctx.contextRunningEventTab.EventContextRunningEvent()
}

// GetCurrentContext 获取当前上下文
func (ctx *ContextBehavior) GetCurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.options.InstanceFace.Iface)
}

// GetConcurrentContext 获取多线程安全的上下文
func (ctx *ContextBehavior) GetConcurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.options.InstanceFace.Iface)
}

// GetInstanceFaceCache 支持重新解释类型
func (ctx *ContextBehavior) GetInstanceFaceCache() iface.Cache {
	return ctx.options.InstanceFace.Cache
}

// CollectGC 收集GC
func (ctx *ContextBehavior) CollectGC(gc GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	ctx.gcList = append(ctx.gcList, gc)
}

// String implements fmt.Stringer
func (ctx *ContextBehavior) String() string {
	ctx.stringerOnce.Do(func() {
		ctx.stringerCache = fmt.Sprintf(`{"name":%q}`, ctx.GetName())
	})
	return ctx.stringerCache
}

func (ctx *ContextBehavior) init(options ContextOptions) {
	ctx.options = options

	if ctx.options.InstanceFace.IsNil() {
		ctx.options.InstanceFace = iface.NewFaceT[Context](ctx)
	}

	if ctx.options.Context == nil {
		ctx.options.Context = context.Background()
	}

	if ctx.options.EntityLib == nil {
		ctx.options.EntityLib = pt.NewEntityLib(pt.NewComponentLib())
	}

	if ctx.options.AddInManager == nil {
		ctx.options.AddInManager = extension.NewRuntimeAddInManager()
	}

	corectx.UnsafeContext(&ctx.ContextBehavior).Init(ctx.options.Context, ctx.options.AutoRecover, ctx.options.ReportError)

	ctx.reflected = reflect.ValueOf(ctx.getInstance())
	ctx.contextRunningEventTab.SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())

	ctx.entityManager.init(ctx.getInstance())

	event.UnsafeEvent(ctx.GetEntityLib().EventEntityLibDeclareEntityPT()).Ctrl().SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())
	event.UnsafeEvent(ctx.GetEntityLib().GetComponentLib().EventComponentLibDeclareComponentPT()).Ctrl().SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())

	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeInstallAddIn()).Ctrl().SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())
	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeUninstallAddIn()).Ctrl().SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())
	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeAddInStateChanged()).Ctrl().SetPanicHandling(ctx.GetAutoRecover(), ctx.GetReportError())

	if ctx.options.RunningEventCB != nil {
		BindEventContextRunningEvent(ctx, HandleEventContextRunningEvent(ctx.options.RunningEventCB))
	}
	BindEventContextRunningEvent(ctx, HandleEventContextRunningEvent(ctx.entityManager.onContextRunningEvent))
}

func (ctx *ContextBehavior) getOptions() *ContextOptions {
	return &ctx.options
}

func (ctx *ContextBehavior) emitEventRunningEvent(runningEvent RunningEvent, args ...any) {
	_EmitEventContextRunningEvent(ctx, ctx.getInstance(), runningEvent, args...)

	switch runningEvent {
	case RunningEvent_Terminated:
		ctx.contextRunningEventTab.SetEnable(false)
		ctx.managed.UnbindAllEventHandles()
	}
}

func (ctx *ContextBehavior) setFrame(frame Frame) {
	ctx.frame = frame
}

func (ctx *ContextBehavior) setCallee(callee async.Callee) {
	ctx.callee = callee
}

func (ctx *ContextBehavior) getScoped() *atomic.Bool {
	return &ctx.scoped
}

func (ctx *ContextBehavior) getInstance() Context {
	return ctx.options.InstanceFace.Iface
}
