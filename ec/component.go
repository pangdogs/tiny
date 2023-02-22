package ec

import (
	"fmt"
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

// Component 组件接口
type Component interface {
	_Component
	_InnerGC
	_InnerGCCollector
	ContextResolver

	// GetID 获取组件ID
	GetID() ID
	// GetName 获取组件名称
	GetName() string
	// GetEntity 获取组件依附的实体
	GetEntity() Entity
	// GetState 获取组件状态
	GetState() ComponentState
	// DestroySelf 销毁自身
	DestroySelf()
	// String 字符串化
	String() string
}

type _Component interface {
	init(name string, entity Entity, inheritor Component, hookAllocator container.Allocator[localevent.Hook])
	setID(id ID)
	setState(state ComponentState)
	getInheritor() Component
	eventComponentDestroySelf() localevent.IEvent
}

// ComponentBehavior 组件行为，需要在开发新组件时，匿名嵌入至组件结构体中
type ComponentBehavior struct {
	id                         ID
	name                       string
	entity                     Entity
	inheritor                  Component
	state                      ComponentState
	_eventComponentDestroySelf localevent.Event
	innerGC                    _ComponentInnerGC
}

// GetID 获取组件ID
func (comp *ComponentBehavior) GetID() ID {
	return comp.id
}

// GetName 获取组件名称
func (comp *ComponentBehavior) GetName() string {
	return comp.name
}

// GetEntity 获取组件依附的实体
func (comp *ComponentBehavior) GetEntity() Entity {
	return comp.entity
}

// GetState 获取组件状态
func (comp *ComponentBehavior) GetState() ComponentState {
	return comp.state
}

// DestroySelf 销毁自身
func (comp *ComponentBehavior) DestroySelf() {
	switch comp.GetState() {
	case ComponentState_Awake, ComponentState_Start, ComponentState_Living:
		emitEventComponentDestroySelf(&comp._eventComponentDestroySelf, comp.inheritor)
	}
}

// String 字符串化
func (comp *ComponentBehavior) String() string {
	var entityID ID
	if entity := comp.GetEntity(); entity != nil {
		entityID = entity.GetID()
	}

	return fmt.Sprintf("[ID:%d Name:%s Entity:%d State:%s]",
		comp.GetID(),
		comp.GetName(),
		entityID,
		comp.GetState())
}

func (comp *ComponentBehavior) init(name string, entity Entity, inheritor Component, hookAllocator container.Allocator[localevent.Hook]) {
	comp.innerGC.Init(comp)
	comp.name = name
	comp.entity = entity
	comp.inheritor = inheritor
	comp._eventComponentDestroySelf.Init(false, nil, localevent.EventRecursion_NotEmit, hookAllocator, &comp.innerGC)
}

func (comp *ComponentBehavior) setID(id ID) {
	comp.id = id
}

func (comp *ComponentBehavior) getContext() util.IfaceCache {
	return comp.entity.getContext()
}

func (comp *ComponentBehavior) setState(state ComponentState) {
	if state <= comp.state {
		return
	}
	comp.state = state
}

func (comp *ComponentBehavior) getInheritor() Component {
	return comp.inheritor
}

func (comp *ComponentBehavior) eventComponentDestroySelf() localevent.IEvent {
	return &comp._eventComponentDestroySelf
}

func (comp *ComponentBehavior) getInnerGC() container.GC {
	return &comp.innerGC
}

func (comp *ComponentBehavior) getInnerGCCollector() container.GCCollector {
	return &comp.innerGC
}
