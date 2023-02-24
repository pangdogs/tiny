package ec

import (
	"kit.golaxy.org/tiny/localevent"
	"kit.golaxy.org/tiny/util"
	"kit.golaxy.org/tiny/util/container"
)

func UnsafeEntity(entity Entity) _UnsafeEntity {
	return _UnsafeEntity{
		Entity: entity,
	}
}

type _UnsafeEntity struct {
	Entity
}

func (ue _UnsafeEntity) Init(opts *EntityOptions) {
	ue.init(opts)
}

func (ue _UnsafeEntity) GetOptions() *EntityOptions {
	return ue.getOptions()
}

func (ue _UnsafeEntity) SetID(id ID) {
	ue.setID(id)
}

func (ue _UnsafeEntity) SetContext(ctx util.IfaceCache) {
	ue.setContext(ctx)
}

func (ue _UnsafeEntity) GetContext() util.IfaceCache {
	return ue.getContext()
}

func (ue _UnsafeEntity) SetGCCollector(gcCollect container.GCCollector) {
	ue.setGCCollector(gcCollect)
}

func (ue _UnsafeEntity) GetGCCollector() container.GCCollector {
	return ue.getGCCollector()
}

func (ue _UnsafeEntity) SetParent(parent Entity) {
	ue.setParent(parent)
}

func (ue _UnsafeEntity) SetState(state EntityState) {
	ue.setState(state)
}

func (ue _UnsafeEntity) EventEntityDestroySelf() localevent.IEvent {
	return ue.eventEntityDestroySelf()
}
