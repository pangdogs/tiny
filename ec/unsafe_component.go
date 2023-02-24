package ec

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

func UnsafeComponent(comp Component) _UnsafeComponent {
	return _UnsafeComponent{
		Component: comp,
	}
}

type _UnsafeComponent struct {
	Component
}

func (uc _UnsafeComponent) Init(name string, entity Entity, inheritor Component, hookAllocator container.Allocator[localevent.Hook], gcCollector container.GCCollector) {
	uc.init(name, entity, inheritor, hookAllocator, gcCollector)
}

func (uc _UnsafeComponent) SetID(id ID) {
	uc.setID(id)
}

func (uc _UnsafeComponent) GetContext() util.IfaceCache {
	return uc.getContext()
}

func (uc _UnsafeComponent) SetState(state ComponentState) {
	uc.setState(state)
}

func (uc _UnsafeComponent) GetInheritor() Component {
	return uc.getInheritor()
}

func (uc _UnsafeComponent) SetGCCollector(gcCollector container.GCCollector) {
	uc.setGCCollector(gcCollector)
}

func (uc _UnsafeComponent) EventComponentDestroySelf() localevent.IEvent {
	return uc.eventComponentDestroySelf()
}
