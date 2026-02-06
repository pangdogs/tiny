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
	"git.golaxy.org/core/utils/uid"
	"git.golaxy.org/tiny/ec/pt"
	"git.golaxy.org/tiny/utils/id"
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

	// Id 获取运行时Id
	Id() uid.Id
	// Name 获取名称
	Name() string
	// Reflected 获取反射值
	Reflected() reflect.Value
	// GenId 生成本地唯一Id
	GenId() id.Id
	// Frame 获取帧
	Frame() Frame
	// EntityManager 获取实体管理器
	EntityManager() EntityManager
	// EntityTree 获取实体树
	EntityTree() EntityTree
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
	idGenerator   int64
	frame         Frame
	entityManager _EntityManager
	callee        async.Callee
	scoped        atomic.Bool
	gcList        []GC
	managed       event.ManagedHandles
	stringerOnce  sync.Once
	stringerCache string

	contextRunningEventTab contextRunningEventTab
}

// Name 获取名称
func (ctx *ContextBehavior) Name() string {
	return ctx.options.Name
}

// Id 获取运行时Id
func (ctx *ContextBehavior) Id() uid.Id {
	return ctx.options.PersistId
}

// Reflected 获取反射值
func (ctx *ContextBehavior) Reflected() reflect.Value {
	return ctx.reflected
}

// GenId 生成本地唯一Id
func (ctx *ContextBehavior) GenId() id.Id {
	if ctx.options.IdGenerator != nil {
		return id.Id(ctx.options.IdGenerator.Add(1))
	}
	ctx.idGenerator++
	return id.Id(ctx.idGenerator)
}

// Frame 获取帧
func (ctx *ContextBehavior) Frame() Frame {
	return ctx.frame
}

// EntityManager 获取实体管理器
func (ctx *ContextBehavior) EntityManager() EntityManager {
	return &ctx.entityManager
}

// EntityTree 获取主实体树
func (ctx *ContextBehavior) EntityTree() EntityTree {
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

// CurrentContext 获取当前上下文
func (ctx *ContextBehavior) CurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.options.InstanceFace.Iface)
}

// ConcurrentContext 获取多线程安全的上下文
func (ctx *ContextBehavior) ConcurrentContext() iface.Cache {
	return iface.Iface2Cache[Context](ctx.options.InstanceFace.Iface)
}

// InstanceFaceCache 支持重新解释类型
func (ctx *ContextBehavior) InstanceFaceCache() iface.Cache {
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
		ctx.stringerCache = fmt.Sprintf(`{"id":%q,"name":%q}`, ctx.Id(), ctx.Name())
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

	if ctx.options.PersistId.IsNil() {
		ctx.options.PersistId = uid.New()
	}

	if ctx.options.EntityLib == nil {
		ctx.options.EntityLib = pt.NewEntityLib(pt.NewComponentLib())
	}

	if ctx.options.AddInManager == nil {
		ctx.options.AddInManager = extension.NewRuntimeAddInManager()
	}

	corectx.UnsafeContext(&ctx.ContextBehavior).Init(ctx.options.Context, ctx.options.AutoRecover, ctx.options.ReportError)

	ctx.reflected = reflect.ValueOf(ctx.getInstance())
	ctx.contextRunningEventTab.SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())

	ctx.entityManager.init(ctx.getInstance())

	event.UnsafeEvent(ctx.EntityLib().EventEntityLibDeclareEntityPT()).Ctrl().SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())
	event.UnsafeEvent(ctx.EntityLib().ComponentLib().EventComponentLibDeclareComponentPT()).Ctrl().SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())

	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeInstallAddIn()).Ctrl().SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())
	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeUninstallAddIn()).Ctrl().SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())
	event.UnsafeEvent(ctx.getAddInManager().EventRuntimeAddInStateChanged()).Ctrl().SetPanicHandling(ctx.AutoRecover(), ctx.ReportError())

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
		ctx.contextRunningEventTab.SetEnabled(false)
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
