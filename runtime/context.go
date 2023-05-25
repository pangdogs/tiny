package runtime

import (
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/internal"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/uid"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// NewContext 创建运行时上下文
func NewContext(options ...ContextOption) Context {
	opts := ContextOptions{}
	WithOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewContext(opts)
}

func UnsafeNewContext(options ContextOptions) Context {
	if !options.CompositeFace.IsNil() {
		options.CompositeFace.Iface.init(&options)
		return options.CompositeFace.Iface
	}

	ctx := &ContextBehavior{}
	ctx.init(&options)

	return ctx.opts.CompositeFace.Iface
}

// Context 运行时上下文接口
type Context interface {
	_Context
	ec.ContextResolver
	container.GCCollector
	internal.Context
	internal.RunningState
	Caller

	// GenPersistId 生成持久化Id
	GenPersistId() uid.Id
	// GetFrame 获取帧
	GetFrame() Frame
	// GetEntityMgr 获取实体管理器
	GetEntityMgr() IEntityMgr
	// GetECTree 获取主EC树
	GetECTree() IECTree
	// GetFaceAnyAllocator 获取FaceAny内存分配器
	GetFaceAnyAllocator() container.Allocator[util.FaceAny]
	// GetHookAllocator 获取Hook内存分配器
	GetHookAllocator() container.Allocator[localevent.Hook]
}

type _Context interface {
	init(opts *ContextOptions)
	getOptions() *ContextOptions
	setFrame(frame Frame)
	setCallee(callee Callee)
	gc()
}

// ContextBehavior 运行时上下文行为，在需要扩展运行时上下文能力时，匿名嵌入至运行时上下文结构体中
type ContextBehavior struct {
	internal.ContextBehavior
	internal.RunningStateBehavior
	opts               ContextOptions
	persistIdGenerator uid.Id
	frame              Frame
	entityMgr          _EntityMgr
	ecTree             ECTree
	callee             Callee
	gcList             []container.GC
}

// GenPersistId 生成持久化Id
func (ctx *ContextBehavior) GenPersistId() uid.Id {
	ctx.persistIdGenerator++
	return ctx.persistIdGenerator
}

// GetFrame 获取帧
func (ctx *ContextBehavior) GetFrame() Frame {
	return ctx.frame
}

// GetEntityMgr 获取实体管理器
func (ctx *ContextBehavior) GetEntityMgr() IEntityMgr {
	return &ctx.entityMgr
}

// GetECTree 获取主EC树
func (ctx *ContextBehavior) GetECTree() IECTree {
	return &ctx.ecTree
}

// GetFaceAnyAllocator 获取FaceAny内存分配器
func (ctx *ContextBehavior) GetFaceAnyAllocator() container.Allocator[util.FaceAny] {
	return ctx.opts.FaceAnyAllocator
}

// GetHookAllocator 获取Hook内存分配器
func (ctx *ContextBehavior) GetHookAllocator() container.Allocator[localevent.Hook] {
	return ctx.opts.HookAllocator
}

// ResolveContext 解析上下文
func (ctx *ContextBehavior) ResolveContext() util.IfaceCache {
	return util.Iface2Cache[Context](ctx.opts.CompositeFace.Iface)
}

// CollectGC 收集GC
func (ctx *ContextBehavior) CollectGC(gc container.GC) {
	if gc == nil || !gc.NeedGC() {
		return
	}

	ctx.gcList = append(ctx.gcList, gc)
}

func (ctx *ContextBehavior) init(opts *ContextOptions) {
	if opts == nil {
		panic("nil opts")
	}

	ctx.opts = *opts

	if ctx.opts.CompositeFace.IsNil() {
		ctx.opts.CompositeFace = util.NewFace[Context](ctx)
	}

	ctx.persistIdGenerator = ctx.opts.PersistIdGenerator

	internal.UnsafeContext(&ctx.ContextBehavior).Init(ctx.opts.Context, ctx.opts.AutoRecover, ctx.opts.ReportError)
	ctx.entityMgr.Init(ctx.getOptions().CompositeFace.Iface)
	ctx.ecTree.init(ctx.opts.CompositeFace.Iface, true)
}

func (ctx *ContextBehavior) getOptions() *ContextOptions {
	return &ctx.opts
}

func (ctx *ContextBehavior) setFrame(frame Frame) {
	ctx.frame = frame
}

func (ctx *ContextBehavior) setCallee(callee Callee) {
	ctx.callee = callee
}
