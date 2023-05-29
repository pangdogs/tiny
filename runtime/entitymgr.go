package runtime

import (
	"errors"
	"fmt"
	"kit.golaxy.org/tiny/ec"
	"kit.golaxy.org/tiny/internal"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/uid"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// IEntityMgr 实体管理器接口
type IEntityMgr interface {
	internal.ContextResolver
	// GetEntity 查询实体
	GetEntity(id uid.Id) (ec.Entity, bool)
	// RangeEntities 遍历所有实体
	RangeEntities(func(entity ec.Entity) bool)
	// ReverseRangeEntities 反向遍历所有实体
	ReverseRangeEntities(func(entity ec.Entity) bool)
	// CountEntities 获取实体数量
	CountEntities() int
	// AddEntity 添加实体
	AddEntity(entity ec.Entity) error
	// RemoveEntity 删除实体
	RemoveEntity(id uid.Id)
	// EventEntityMgrAddEntity 事件：实体管理器添加实体
	EventEntityMgrAddEntity() localevent.IEvent
	// EventEntityMgrRemovingEntity 事件：实体管理器开始删除实体
	EventEntityMgrRemovingEntity() localevent.IEvent
	// EventEntityMgrRemoveEntity 事件：实体管理器删除实体
	EventEntityMgrRemoveEntity() localevent.IEvent
	// EventEntityMgrEntityAddComponents 事件：实体管理器中的实体添加组件
	EventEntityMgrEntityAddComponents() localevent.IEvent
	// EventEntityMgrEntityRemoveComponent 事件：实体管理器中的实体删除组件
	EventEntityMgrEntityRemoveComponent() localevent.IEvent
}

type _EntityInfo struct {
	Element *container.Element[util.FaceAny]
	Hooks   [2]localevent.Hook
}

type _EntityMgr struct {
	ctx                                 Context
	entityMap                           map[uid.Id]_EntityInfo
	entityList                          container.List[util.FaceAny]
	eventEntityMgrAddEntity             localevent.Event
	eventEntityMgrRemovingEntity        localevent.Event
	eventEntityMgrRemoveEntity          localevent.Event
	eventEntityMgrEntityAddComponents   localevent.Event
	eventEntityMgrEntityRemoveComponent localevent.Event
}

func (entityMgr *_EntityMgr) Init(ctx Context) {
	if ctx == nil {
		panic("nil ctx")
	}

	entityMgr.ctx = ctx
	entityMgr.entityList.Init(ctx.GetFaceAnyAllocator(), ctx)
	entityMgr.entityMap = map[uid.Id]_EntityInfo{}

	entityMgr.eventEntityMgrAddEntity.Init(ctx.GetAutoRecover(), ctx.GetReportError(), localevent.EventRecursion_Allow, ctx.GetHookAllocator(), ctx)
	entityMgr.eventEntityMgrRemovingEntity.Init(ctx.GetAutoRecover(), ctx.GetReportError(), localevent.EventRecursion_Allow, ctx.GetHookAllocator(), ctx)
	entityMgr.eventEntityMgrRemoveEntity.Init(ctx.GetAutoRecover(), ctx.GetReportError(), localevent.EventRecursion_Allow, ctx.GetHookAllocator(), ctx)
	entityMgr.eventEntityMgrEntityAddComponents.Init(ctx.GetAutoRecover(), ctx.GetReportError(), localevent.EventRecursion_Allow, ctx.GetHookAllocator(), ctx)
	entityMgr.eventEntityMgrEntityRemoveComponent.Init(ctx.GetAutoRecover(), ctx.GetReportError(), localevent.EventRecursion_Allow, ctx.GetHookAllocator(), ctx)
}

// ResolveContext 解析上下文
func (entityMgr *_EntityMgr) ResolveContext() util.IfaceCache {
	return entityMgr.ctx.ResolveContext()
}

// GetEntity 查询实体
func (entityMgr *_EntityMgr) GetEntity(id uid.Id) (ec.Entity, bool) {
	e, ok := entityMgr.entityMap[id]
	if !ok {
		return nil, false
	}

	if e.Element.Escaped() {
		return nil, false
	}

	return util.Cache2Iface[ec.Entity](e.Element.Value.Cache), true
}

// RangeEntities 遍历所有实体
func (entityMgr *_EntityMgr) RangeEntities(fun func(entity ec.Entity) bool) {
	if fun == nil {
		return
	}

	entityMgr.entityList.Traversal(func(e *container.Element[util.FaceAny]) bool {
		return fun(util.Cache2Iface[ec.Entity](e.Value.Cache))
	})
}

// ReverseRangeEntities 反向遍历所有实体
func (entityMgr *_EntityMgr) ReverseRangeEntities(fun func(entity ec.Entity) bool) {
	if fun == nil {
		return
	}

	entityMgr.entityList.ReverseTraversal(func(e *container.Element[util.FaceAny]) bool {
		return fun(util.Cache2Iface[ec.Entity](e.Value.Cache))
	})
}

// CountEntities 获取实体数量
func (entityMgr *_EntityMgr) CountEntities() int {
	return entityMgr.entityList.Len()
}

// AddEntity 添加实体
func (entityMgr *_EntityMgr) AddEntity(entity ec.Entity) error {
	if entity == nil {
		return errors.New("nil entity")
	}

	if entity.GetState() != ec.EntityState_Birth {
		return errors.New("entity state not birth is invalid")
	}

	if entity.GetId() != util.Zero[uid.Id]() {
		if _, ok := entityMgr.entityMap[entity.GetId()]; ok {
			return fmt.Errorf("entity id already existed")
		}
	}

	ctx := entityMgr.ctx
	_entity := ec.UnsafeEntity(entity)

	if _entity.GetId() == util.Zero[uid.Id]() {
		_entity.SetId(ctx.GenPersistId())
	}
	_entity.SetContext(util.Iface2Cache[Context](ctx))

	_entity.RangeComponents(func(comp ec.Component) bool {
		_comp := ec.UnsafeComponent(comp)

		if _comp.GetId() == util.Zero[uid.Id]() {
			_comp.SetId(ctx.GenPersistId())
		}

		return true
	})

	entityInfo := _EntityInfo{}
	entityInfo.Element = entityMgr.entityList.PushBack(util.NewFacePair[any](entity, entity))
	entityInfo.Hooks[0] = localevent.BindEvent[ec.EventCompMgrAddComponents](entity.EventCompMgrAddComponents(), entityMgr)
	entityInfo.Hooks[1] = localevent.BindEvent[ec.EventCompMgrRemoveComponent](entity.EventCompMgrRemoveComponent(), entityMgr)

	entityMgr.entityMap[entity.GetId()] = entityInfo

	_entity.SetState(ec.EntityState_Entry)

	if _entity.GetGCCollector() == nil {
		_entity.SetGCCollector(ctx)
	}

	emitEventEntityMgrAddEntity(&entityMgr.eventEntityMgrAddEntity, entityMgr, entity)

	return nil
}

// RemoveEntity 删除实体
func (entityMgr *_EntityMgr) RemoveEntity(id uid.Id) {
	entityInfo, ok := entityMgr.entityMap[id]
	if !ok {
		return
	}

	entity := ec.UnsafeEntity(util.Cache2Iface[ec.Entity](entityInfo.Element.Value.Cache))
	if entity.GetState() >= ec.EntityState_Leave {
		return
	}

	entity.SetState(ec.EntityState_Leave)

	emitEventEntityMgrRemovingEntity(&entityMgr.eventEntityMgrRemovingEntity, entityMgr, entity.Entity)

	delete(entityMgr.entityMap, id)
	entityInfo.Element.Escape()

	for i := range entityInfo.Hooks {
		entityInfo.Hooks[i].Unbind()
	}

	emitEventEntityMgrRemoveEntity(&entityMgr.eventEntityMgrRemoveEntity, entityMgr, entity.Entity)
}

// EventEntityMgrAddEntity 事件：实体管理器添加实体
func (entityMgr *_EntityMgr) EventEntityMgrAddEntity() localevent.IEvent {
	return &entityMgr.eventEntityMgrAddEntity
}

// EventEntityMgrRemovingEntity 事件：实体管理器开始删除实体
func (entityMgr *_EntityMgr) EventEntityMgrRemovingEntity() localevent.IEvent {
	return &entityMgr.eventEntityMgrRemovingEntity
}

// EventEntityMgrRemoveEntity 事件：实体管理器删除实体
func (entityMgr *_EntityMgr) EventEntityMgrRemoveEntity() localevent.IEvent {
	return &entityMgr.eventEntityMgrRemoveEntity
}

// EventEntityMgrEntityAddComponents 事件：实体管理器中的实体添加组件
func (entityMgr *_EntityMgr) EventEntityMgrEntityAddComponents() localevent.IEvent {
	return &entityMgr.eventEntityMgrEntityAddComponents
}

// EventEntityMgrEntityRemoveComponent 事件：实体管理器中的实体删除组件
func (entityMgr *_EntityMgr) EventEntityMgrEntityRemoveComponent() localevent.IEvent {
	return &entityMgr.eventEntityMgrEntityRemoveComponent
}

func (entityMgr *_EntityMgr) OnCompMgrAddComponents(entity ec.Entity, components []ec.Component) {
	for i := range components {
		_comp := ec.UnsafeComponent(components[i])

		if _comp.GetId() == util.Zero[uid.Id]() {
			_comp.SetId(entityMgr.ctx.GenPersistId())
		}
	}

	emitEventEntityMgrEntityAddComponents(&entityMgr.eventEntityMgrEntityAddComponents, entityMgr, entity, components)
}

func (entityMgr *_EntityMgr) OnCompMgrRemoveComponent(entity ec.Entity, component ec.Component) {
	emitEventEntityMgrEntityRemoveComponent(&entityMgr.eventEntityMgrEntityRemoveComponent, entityMgr, entity, component)
}
