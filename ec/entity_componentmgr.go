package ec

import (
	"fmt"
	"git.golaxy.org/tiny/event"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/generic"
	"git.golaxy.org/tiny/utils/iface"
	"git.golaxy.org/tiny/utils/uid"
	"golang.org/x/exp/slices"
)

// iComponentMgr 组件管理器接口
type iComponentMgr interface {
	// GetComponent 使用名称查询组件，同个名称指向多个组件时，返回首个组件
	GetComponent(name string) Component
	// GetComponentById 使用组件Id查询组件
	GetComponentById(id uid.Id) Component
	// ContainsComponent 组件是否存在
	ContainsComponent(name string) bool
	// ContainsComponentById 使用组件Id检测组件是否存在
	ContainsComponentById(id uid.Id) bool
	// RangeComponents 遍历所有组件
	RangeComponents(fun generic.Func1[Component, bool])
	// ReversedRangeComponents 反向遍历所有组件
	ReversedRangeComponents(fun generic.Func1[Component, bool])
	// FilterComponents 过滤并获取组件
	FilterComponents(fun generic.Func1[Component, bool]) []Component
	// GetComponents 获取所有组件
	GetComponents() []Component
	// CountComponents 统计所有组件数量
	CountComponents() int
	// AddComponent 添加组件，因为同个名称可以指向多个组件，所以名称指向的组件已存在时，不会返回错误
	AddComponent(name string, components ...Component) error
	// RemoveComponent 使用名称删除组件，将会删除同个名称指向的多个组件
	RemoveComponent(name string)
	// RemoveComponentById 使用组件Id删除组件
	RemoveComponentById(id uid.Id)

	iAutoEventComponentMgrAddComponents        // 事件：实体的组件管理器添加组件
	iAutoEventComponentMgrRemoveComponent      // 事件：实体的组件管理器删除组件
	iAutoEventComponentMgrFirstAccessComponent // 事件：实体的组件管理器首次访问组件
}

// GetComponent 使用名称查询组件，同个名称指向多个组件时，返回首个组件
func (entity *EntityBehavior) GetComponent(name string) Component {
	if n, ok := entity.getComponentNode(name); ok {
		return entity.accessComponent(iface.Cache2Iface[Component](n.V.Cache))
	}
	return nil
}

// GetComponentById 使用组件Id查询组件
func (entity *EntityBehavior) GetComponentById(id uid.Id) Component {
	if n, ok := entity.getComponentNodeById(id); ok {
		return entity.accessComponent(iface.Cache2Iface[Component](n.V.Cache))
	}
	return nil
}

// ContainsComponent 组件是否存在
func (entity *EntityBehavior) ContainsComponent(name string) bool {
	_, ok := entity.getComponentNode(name)
	return ok
}

// ContainsComponentById 使用组件Id检测组件是否存在
func (entity *EntityBehavior) ContainsComponentById(id uid.Id) bool {
	_, ok := entity.getComponentNodeById(id)
	return ok
}

// RangeComponents 遍历所有组件
func (entity *EntityBehavior) RangeComponents(fun generic.Func1[Component, bool]) {
	entity.componentList.Traversal(func(n *generic.Node[iface.FaceAny]) bool {
		comp := entity.accessComponent(iface.Cache2Iface[Component](n.V.Cache))
		if comp == nil {
			return true
		}
		return fun.Exec(comp)
	})
}

// ReversedRangeComponents 反向遍历所有组件
func (entity *EntityBehavior) ReversedRangeComponents(fun generic.Func1[Component, bool]) {
	entity.componentList.ReversedTraversal(func(n *generic.Node[iface.FaceAny]) bool {
		comp := entity.accessComponent(iface.Cache2Iface[Component](n.V.Cache))
		if comp == nil {
			return true
		}
		return fun.Exec(comp)
	})
}

// FilterComponents 过滤并获取组件
func (entity *EntityBehavior) FilterComponents(fun generic.Func1[Component, bool]) []Component {
	var components []Component

	entity.componentList.Traversal(func(n *generic.Node[iface.FaceAny]) bool {
		comp := iface.Cache2Iface[Component](n.V.Cache)
		if fun.Exec(comp) {
			components = append(components, comp)
		}
		return true
	})

	for i := range components {
		if entity.accessComponent(components[i]) == nil {
			components[i] = nil
		}
	}

	components = slices.DeleteFunc(components, func(comp Component) bool {
		return comp == nil
	})

	return components
}

// GetComponents 获取所有组件
func (entity *EntityBehavior) GetComponents() []Component {
	components := make([]Component, 0, entity.componentList.Len())

	entity.componentList.Traversal(func(n *generic.Node[iface.FaceAny]) bool {
		components = append(components, iface.Cache2Iface[Component](n.V.Cache))
		return true
	})

	for i := range components {
		if entity.accessComponent(components[i]) == nil {
			components[i] = nil
		}
	}

	components = slices.DeleteFunc(components, func(comp Component) bool {
		return comp == nil
	})

	return components
}

// CountComponents 统计所有组件数量
func (entity *EntityBehavior) CountComponents() int {
	return entity.componentList.Len()
}

// AddComponent 添加组件，因为同个名称可以指向多个组件，所以名称指向的组件已存在时，不会返回错误
func (entity *EntityBehavior) AddComponent(name string, components ...Component) error {
	if len(components) <= 0 {
		return fmt.Errorf("%w: %w: components is empty", ErrEC, exception.ErrArgs)
	}

	for i := range components {
		comp := components[i]

		if comp == nil {
			return fmt.Errorf("%w: %w: component is nil", ErrEC, exception.ErrArgs)
		}

		if comp.GetState() != ComponentState_Birth {
			return fmt.Errorf("%w: invalid component state %q", ErrEC, comp.GetState())
		}
	}

	for i := range components {
		if err := entity.addComponent(name, components[i]); err != nil {
			return err
		}
	}

	_EmitEventComponentMgrAddComponents(entity, entity.opts.CompositeFace.Iface, components)
	return nil
}

// RemoveComponent 使用名称删除组件，将会删除同个名称指向的多个组件
func (entity *EntityBehavior) RemoveComponent(name string) {
	n, ok := entity.getComponentNode(name)
	if !ok {
		return
	}

	entity.componentList.TraversalAt(func(other *generic.Node[iface.FaceAny]) bool {
		comp := iface.Cache2Iface[Component](other.V.Cache)
		if comp.GetName() != name {
			return false
		}

		if comp.GetState() > ComponentState_Alive {
			return true
		}
		comp.setState(ComponentState_Detach)

		_EmitEventComponentMgrRemoveComponent(entity, entity.opts.CompositeFace.Iface, comp)

		other.Escape()
		return true
	}, n)
}

// RemoveComponentById 使用组件Id删除组件
func (entity *EntityBehavior) RemoveComponentById(id uid.Id) {
	n, ok := entity.getComponentNodeById(id)
	if !ok {
		return
	}

	comp := iface.Cache2Iface[Component](n.V.Cache)

	if comp.GetState() > ComponentState_Alive {
		return
	}
	comp.setState(ComponentState_Detach)

	_EmitEventComponentMgrRemoveComponent(entity, entity.opts.CompositeFace.Iface, comp)

	n.Escape()
}

// EventComponentMgrAddComponents 事件：实体的组件管理器添加组件
func (entity *EntityBehavior) EventComponentMgrAddComponents() event.IEvent {
	return &entity.eventComponentMgrAddComponents
}

// EventComponentMgrRemoveComponent 事件：实体的组件管理器删除组件
func (entity *EntityBehavior) EventComponentMgrRemoveComponent() event.IEvent {
	return &entity.eventComponentMgrRemoveComponent
}

// EventComponentMgrFirstAccessComponent 事件：实体的组件管理器首次访问组件
func (entity *EntityBehavior) EventComponentMgrFirstAccessComponent() event.IEvent {
	return &entity.eventComponentMgrFirstAccessComponent
}

func (entity *EntityBehavior) addComponent(name string, component Component) error {
	component.init(name, entity.opts.CompositeFace.Iface, component)

	face := iface.MakeFaceAny(component)

	if n, ok := entity.getComponentNode(name); ok {
		entity.componentList.TraversalAt(func(other *generic.Node[iface.FaceAny]) bool {
			if iface.Cache2Iface[Component](other.V.Cache).GetName() == name {
				n = other
				return true
			}
			return false
		}, n)

		n = entity.componentList.InsertAfter(face, n)

	} else {
		n = entity.componentList.PushBack(face)
	}

	component.setState(ComponentState_Attach)

	return nil
}

func (entity *EntityBehavior) getComponentNode(name string) (*generic.Node[iface.FaceAny], bool) {
	var n *generic.Node[iface.FaceAny]

	entity.componentList.Traversal(func(other *generic.Node[iface.FaceAny]) bool {
		if iface.Cache2Iface[Component](other.V.Cache).GetName() == name {
			n = other
			return false
		}
		return true
	})

	return n, n != nil
}

func (entity *EntityBehavior) getComponentNodeById(id uid.Id) (*generic.Node[iface.FaceAny], bool) {
	var n *generic.Node[iface.FaceAny]

	entity.componentList.Traversal(func(other *generic.Node[iface.FaceAny]) bool {
		if iface.Cache2Iface[Component](other.V.Cache).GetId() == id {
			n = other
			return false
		}
		return true
	})

	return n, n != nil
}

func (entity *EntityBehavior) accessComponent(comp Component) Component {
	if entity.opts.AwakeOnFirstAccess && comp.GetState() == ComponentState_Attach {
		switch entity.GetState() {
		case EntityState_Awake, EntityState_Start, EntityState_Alive:
			_EmitEventComponentMgrFirstAccessComponent(entity, entity.opts.CompositeFace.Iface, comp)
		}
	}

	if comp.GetState() >= ComponentState_Death {
		return nil
	}

	return comp
}
