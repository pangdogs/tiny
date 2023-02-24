package ec

import (
	"fmt"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// NewEntity 创建实体
func NewEntity(options ...EntityOption) Entity {
	opts := EntityOptions{}
	WithEntityOption{}.Default()(&opts)

	for i := range options {
		options[i](&opts)
	}

	return UnsafeNewEntity(opts)
}

func UnsafeNewEntity(options EntityOptions) Entity {
	if !options.Inheritor.IsNil() {
		options.Inheritor.Iface.init(&options)
		return options.Inheritor.Iface
	}

	e := &EntityBehavior{}
	e.init(&options)

	return e.opts.Inheritor.Iface
}

// Entity 实体接口
type Entity interface {
	_Entity
	_ComponentMgr
	ContextResolver

	// GetID 获取实体ID
	GetID() ID
	// GetParent 获取在运行时上下文的主EC树上的父实体
	GetParent() (Entity, bool)
	// GetState 获取实体状态
	GetState() EntityState
	// DestroySelf 销毁自身
	DestroySelf()
	// String 字符串化
	String() string
}

type _Entity interface {
	init(opts *EntityOptions)
	getOptions() *EntityOptions
	setID(id ID)
	setContext(ctx util.IfaceCache)
	setGCCollector(gcCollector container.GCCollector)
	getGCCollector() container.GCCollector
	setParent(parent Entity)
	setState(state EntityState)
	eventEntityDestroySelf() localevent.IEvent
}

// EntityBehavior 实体行为，在需要扩展实体能力时，匿名嵌入至实体结构体中
type EntityBehavior struct {
	id                          ID
	opts                        EntityOptions
	context                     util.IfaceCache
	parent                      Entity
	componentList               container.List[util.FaceAny]
	state                       EntityState
	_eventEntityDestroySelf     localevent.Event
	eventCompMgrAddComponents   localevent.Event
	eventCompMgrRemoveComponent localevent.Event
}

// GetID 获取实体ID
func (entity *EntityBehavior) GetID() ID {
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
	case EntityState_Init, EntityState_Start, EntityState_Living:
		emitEventEntityDestroySelf(&entity._eventEntityDestroySelf, entity.opts.Inheritor.Iface)
	}
}

// String 字符串化
func (entity *EntityBehavior) String() string {
	var parentID ID
	if parent, ok := entity.GetParent(); ok {
		parentID = parent.GetID()
	}

	return fmt.Sprintf("[ID:%d Parent:%d State:%s]",
		entity.GetID(),
		parentID,
		entity.GetState())
}

func (entity *EntityBehavior) init(opts *EntityOptions) {
	if opts == nil {
		panic("nil opts")
	}

	entity.opts = *opts

	if entity.opts.Inheritor.IsNil() {
		entity.opts.Inheritor = util.NewFace[Entity](entity)
	}

	entity.componentList.Init(entity.opts.FaceAnyAllocator, entity.opts.GCCollector)

	entity._eventEntityDestroySelf.Init(false, nil, localevent.EventRecursion_NotEmit, entity.opts.HookAllocator, entity.opts.GCCollector)
	entity.eventCompMgrAddComponents.Init(false, nil, localevent.EventRecursion_Allow, entity.opts.HookAllocator, entity.opts.GCCollector)
	entity.eventCompMgrRemoveComponent.Init(false, nil, localevent.EventRecursion_Allow, entity.opts.HookAllocator, entity.opts.GCCollector)
}

func (entity *EntityBehavior) getOptions() *EntityOptions {
	return &entity.opts
}

func (entity *EntityBehavior) setID(id ID) {
	entity.id = id
}

func (entity *EntityBehavior) setContext(ctx util.IfaceCache) {
	entity.context = ctx
}

func (entity *EntityBehavior) getContext() util.IfaceCache {
	return entity.context
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
