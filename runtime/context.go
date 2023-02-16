package runtime

import (
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/internal"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// NewContext 创建运行时上下文
func NewContext(options ...ContextOption) Context {
	opts := ContextOptions{}
	WithContextOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewContext(opts)
}

func UnsafeNewContext(options ContextOptions) Context {
	if !options.Inheritor.IsNil() {
		options.Inheritor.Iface.init(&options)
		return options.Inheritor.Iface
	}

	ctx := &ContextBehavior{}
	ctx.init(&options)

	return ctx.opts.Inheritor.Iface
}

// Context 运行时上下文接口
type Context interface {
	_Context
	_InnerGC
	container.GCCollector
	internal.Context
	internal.RunningMark
	_Call

	// GenPersistID 生成持久化ID
	GenPersistID() ec.ID
	// GetFrame 获取帧
	GetFrame() Frame
	// GetEntityMgr 获取实体管理器
	GetEntityMgr() IEntityMgr
	// GetECTree 获取主EC树
	GetECTree() IECTree
	// GetFaceCache 获取Face缓存
	GetFaceCache() container.Cache[util.FaceAny]
	// GetHookCache 获取Hook缓存
	GetHookCache() container.Cache[localevent.Hook]
}

type _Context interface {
	init(opts *ContextOptions)
	getOptions() *ContextOptions
	setFrame(frame Frame)
}

// ContextBehavior 运行时上下文行为，在需要扩展运行时上下文能力时，匿名嵌入至运行时上下文结构体中
type ContextBehavior struct {
	internal.ContextBehavior
	internal.RunningMarkBehavior
	opts               ContextOptions
	persistIDGenerator ec.ID
	frame              Frame
	entityMgr          _EntityMgr
	ecTree             ECTree
	callee             internal.Callee
	gcList             []container.GC
	innerGC            _ContextInnerGC
}

// GenPersistID 生成持久化ID
func (ctx *ContextBehavior) GenPersistID() ec.ID {
	ctx.persistIDGenerator++
	return ctx.persistIDGenerator
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

// GetFaceCache 获取Face缓存
func (ctx *ContextBehavior) GetFaceCache() container.Cache[util.FaceAny] {
	return ctx.opts.FaceCache
}

// GetHookCache 获取Hook缓存
func (ctx *ContextBehavior) GetHookCache() container.Cache[localevent.Hook] {
	return ctx.opts.HookCache
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

	if ctx.opts.Inheritor.IsNil() {
		ctx.opts.Inheritor = util.NewFace[Context](ctx)
	}

	ctx.persistIDGenerator = ctx.opts.PersistIDGenerator

	ctx.innerGC.Init(ctx)

	internal.UnsafeContext(&ctx.ContextBehavior).Init(ctx.opts.Context, ctx.opts.AutoRecover, ctx.opts.ReportError)
	ctx.entityMgr.Init(ctx.getOptions().Inheritor.Iface)
	ctx.ecTree.init(ctx.opts.Inheritor.Iface, true)
}

func (ctx *ContextBehavior) getOptions() *ContextOptions {
	return &ctx.opts
}

func (ctx *ContextBehavior) setFrame(frame Frame) {
	ctx.frame = frame
}

func (ctx *ContextBehavior) getInnerGC() container.GC {
	return &ctx.innerGC
}
