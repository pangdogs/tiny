package ec

import (
	"fmt"
	"kit.golaxy.org/tiny/internal"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/uid"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// NewEntity 创建实体
func NewEntity(options ...EntityOption) Entity {
	opts := EntityOptions{}
	WithOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewEntity(opts)
}

func UnsafeNewEntity(options EntityOptions) Entity {
	if !options.CompositeFace.IsNil() {
		options.CompositeFace.Iface.init(&options)
		return options.CompositeFace.Iface
	}

	e := &EntityBehavior{}
	e.init(&options)

	return e.opts.CompositeFace.Iface
}

// Entity 实体接口
type Entity interface {
	_Entity
	_ComponentMgr
	internal.ContextResolver
	fmt.Stringer

	// GetId 获取实体Id
	GetId() uid.Id
	// GetParent 获取在运行时上下文的主EC树上的父实体
	GetParent() (Entity, bool)
	// GetState 获取实体状态
	GetState() EntityState
	// DestroySelf 销毁自身
	DestroySelf()
}

type _Entity interface {
	init(opts *EntityOptions)
	getOptions() *EntityOptions
	setId(id uid.Id)
	setContext(ctx util.IfaceCache)
	setGCCollector(gcCollector container.GCCollector)
	getGCCollector() container.GCCollector
	setParent(parent Entity)
	setState(state EntityState)
	eventEntityDestroySelf() localevent.IEvent
}

// EntityBehavior 实体行为，在需要扩展实体能力时，匿名嵌入至实体结构体中
type EntityBehavior struct {
	id                          uid.Id
	opts                        EntityOptions
	context                     util.IfaceCache
	parent                      Entity
	componentList               container.List[util.FaceAny]
	state                       EntityState
	_eventEntityDestroySelf     localevent.Event
	eventCompMgrAddComponents   localevent.Event
	eventCompMgrRemoveComponent localevent.Event
}

// GetId 获取实体Id
func (entity *EntityBehavior) GetId() uid.Id {
	return entity.id
}

// GetParent 获取在运行时上下文的主EC树上的父实体
func (entity *EntityBehavior) GetParent() (Entity, bool) {
	return entity.parent, entity.parent != nil
}

// GetState 获取实体状态
func (entity *EntityBehavior) GetState() EntityState {
	return entity.state
}

// DestroySelf 销毁自身
func (entity *EntityBehavior) DestroySelf() {
	switch entity.GetState() {
	case EntityState_Init, EntityState_Inited, EntityState_Living:
		emitEventEntityDestroySelf(&entity._eventEntityDestroySelf, entity.opts.CompositeFace.Iface)
	}
}

// ResolveContext 解析上下文
func (entity *EntityBehavior) ResolveContext() util.IfaceCache {
	return entity.context
}

// String 字符串化
func (entity *EntityBehavior) String() string {
	var parentId uid.Id
	if parent, ok := entity.GetParent(); ok {
		parentId = parent.GetId()
	}

	return fmt.Sprintf("{Id:%d Parent:%d State:%s}", entity.GetId(), parentId, entity.GetState())
}

func (entity *EntityBehavior) init(opts *EntityOptions) {
	if opts == nil {
		panic("nil opts")
	}

	entity.opts = *opts

	if entity.opts.CompositeFace.IsNil() {
		entity.opts.CompositeFace = util.NewFace[Entity](entity)
	}

	entity.componentList.Init(entity.opts.FaceAnyAllocator, entity.opts.GCCollector)

	entity._eventEntityDestroySelf.Init(false, nil, localevent.EventRecursion_NotEmit, entity.opts.HookAllocator, entity.opts.GCCollector)
	entity.eventCompMgrAddComponents.Init(false, nil, localevent.EventRecursion_Allow, entity.opts.HookAllocator, entity.opts.GCCollector)
	entity.eventCompMgrRemoveComponent.Init(false, nil, localevent.EventRecursion_Allow, entity.opts.HookAllocator, entity.opts.GCCollector)
}

func (entity *EntityBehavior) getOptions() *EntityOptions {
	return &entity.opts
}

func (entity *EntityBehavior) setId(id uid.Id) {
	entity.id = id
}

func (entity *EntityBehavior) setContext(ctx util.IfaceCache) {
	entity.context = ctx
}

func (entity *EntityBehavior) setGCCollector(gcCollector container.GCCollector) {
	if entity.opts.GCCollector == gcCollector {
		return
	}

	entity.opts.GCCollector = gcCollector

	entity.componentList.SetGCCollector(gcCollector)
	entity.componentList.Traversal(func(e *container.Element[util.FaceAny]) bool {
		comp := util.Cache2Iface[Component](e.Value.Cache)
		comp.setGCCollector(gcCollector)
		return true
	})

	localevent.UnsafeEvent(&entity._eventEntityDestroySelf).SetGCCollector(gcCollector)
	localevent.UnsafeEvent(&entity.eventCompMgrAddComponents).SetGCCollector(gcCollector)
	localevent.UnsafeEvent(&entity.eventCompMgrRemoveComponent).SetGCCollector(gcCollector)
}

func (entity *EntityBehavior) getGCCollector() container.GCCollector {
	return entity.opts.GCCollector
}

func (entity *EntityBehavior) setParent(parent Entity) {
	entity.parent = parent
}

func (entity *EntityBehavior) setState(state EntityState) {
	if state <= entity.state {
		return
	}
	entity.state = state
}

func (entity *EntityBehavior) eventEntityDestroySelf() localevent.IEvent {
	return &entity._eventEntityDestroySelf
}
