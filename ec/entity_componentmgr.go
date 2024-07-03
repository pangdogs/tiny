package ec

import (
	"fmt"
	"git.golaxy.org/tiny/event"
	"git.golaxy.org/tiny/utils/exception"
	"git.golaxy.org/tiny/utils/generic"
	"git.golaxy.org/tiny/utils/uid"
	"golang.org/x/exp/slices"
)

// iComponentMgr 组件管理器接口
type iComponentMgr interface {
	使用名称查询组件，组件同名时，返回首个组件
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
	// AddComponent 添加组件，允许组件同名
	AddComponent(name string, components ...Component) error
	// AddFixedComponent 添加固定组件，不允许组件同名
	AddFixedComponent(name string, components ...Component) error
	// RemoveComponent 使用名称删除组件，同名组件均会删除
	RemoveComponent(name string)
	// RemoveComponentById 使用组件Id删除组件
	RemoveComponentById(id uid.Id)

	iAutoEventComponentMgrAddComponents        // 事件：实体的组件管理器添加组件
	iAutoEventComponentMgrRemoveComponent      // 事件：实体的组件管理器删除组件
	iAutoEventComponentMgrFirstAccessComponent // 事件：实体的组件管理器首次访问组件
}

// GetComponent 使用名称查询组件，组件同名时，返回首个组件
func (entity *EntityBehavior) GetComponent(name string) Component {
	comp, ok := entity.getFixedComponentNode(name)
	if ok {
		return entity.accessComponent(comp)
	}
	node, ok := entity.getMutableComponentNode(name)
	if ok {
		return entity.accessComponent(node.V)
	}
	return nil
}

// GetComponentById 使用组件Id查询组件
func (entity *EntityBehavior) GetComponentById(id uid.Id) Component {
	comp, ok := entity.getFixedComponentNodeById(id)
	if ok {
		return entity.accessComponent(comp)
	}
	node, ok := entity.getMutableComponentNodeById(id)
	if ok {
		return entity.accessComponent(node.V)
	}
	return nil
}

// ContainsComponent 组件是否存在
func (entity *EntityBehavior) ContainsComponent(name string) bool {
	_, ok := entity.getFixedComponentNode(name)
	if ok {
		return true
	}
	_, ok = entity.getMutableComponentNode(name)
	return ok
}

// ContainsComponentById 使用组件Id检测组件是否存在
func (entity *EntityBehavior) ContainsComponentById(id uid.Id) bool {
	_, ok := entity.getMutableComponentNodeById(id)
	if ok {
		return true
	}
	_, ok = entity.getMutableComponentNodeById(id)
	return ok
}

// RangeComponents 遍历所有组件
func (entity *EntityBehavior) RangeComponents(fun generic.Func1[Component, bool]) {
	_continue := true

	entity.fixedComponentList.Range(func(_ string, comp Component) bool {
		comp = entity.accessComponent(comp)
		if comp == nil {
			return true
		}
		_continue = fun.Exec(comp)
		return _continue
	})

	if !_continue {
		return
	}

	entity.mutableComponentList.Traversal(func(node *generic.Node[Component]) bool {
		comp := entity.accessComponent(node.V)
		if comp == nil {
			return true
		}
		return fun.Exec(comp)
	})
}

// ReversedRangeComponents 反向遍历所有组件
func (entity *EntityBehavior) ReversedRangeComponents(fun generic.Func1[Component, bool]) {
	_continue := true

	entity.mutableComponentList.ReversedTraversal(func(node *generic.Node[Component]) bool {
		comp := entity.accessComponent(node.V)
		if comp == nil {
			return true
		}
		_continue = fun.Exec(comp)
		return _continue
	})

	if !_continue {
		return
	}

	entity.fixedComponentList.ReversedRange(func(_ string, comp Component) bool {
		comp = entity.accessComponent(comp)
		if comp == nil {
			return true
		}
		return fun.Exec(comp)
	})
}

// FilterComponents 过滤并获取组件
func (entity *EntityBehavior) FilterComponents(fun generic.Func1[Component, bool]) []Component {
	var comps []Component

	entity.fixedComponentList.Each(func(_ string, comp Component) {
		if fun.Exec(comp) {
			comps = append(comps, comp)
		}
	})

	entity.mutableComponentList.Traversal(func(node *generic.Node[Component]) bool {
		comp := node.V
		if fun.Exec(comp) {
			comps = append(comps, comp)
		}
		return true
	})

	for i := range comps {
		if entity.accessComponent(comps[i]) == nil {
			comps[i] = nil
		}
	}

	comps = slices.DeleteFunc(comps, func(comp Component) bool {
		return comp == nil
	})

	return comps
}

// GetComponents 获取所有组件
func (entity *EntityBehavior) GetComponents() []Component {
	comps := make([]Component, 0, entity.CountComponents())

	entity.fixedComponentList.Each(func(_ string, comp Component) {
		comps = append(comps, comp)
	})

	entity.mutableComponentList.Traversal(func(node *generic.Node[Component]) bool {
		comps = append(comps, node.V)
		return true
	})

	for i := range comps {
		if entity.accessComponent(comps[i]) == nil {
			comps[i] = nil
		}
	}

	comps = slices.DeleteFunc(comps, func(comp Component) bool {
		return comp == nil
	})

	return comps
}

// CountComponents 统计所有组件数量
func (entity *EntityBehavior) CountComponents() int {
	return entity.fixedComponentList.Len() + entity.mutableComponentList.Len()
}

// AddComponent 添加组件，允许组件同名
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
		entity.addMutableComponent(name, components[i])
	}

	_EmitEventComponentMgrAddComponents(entity, entity.opts.CompositeFace.Iface, components)
	return nil
}

// AddFixedComponent 添加固定组件，不允许组件同名
func (entity *EntityBehavior) AddFixedComponent(name string, components ...Component) error {
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

		if entity.fixedComponentList.Exist(comp.GetName()) {
			return fmt.Errorf("%w: fixed component %q already exists", ErrEC, comp.GetState())
		}
	}

	for i := range components {
		entity.addFixedComponent(name, components[i])
	}

	_EmitEventComponentMgrAddComponents(entity, entity.opts.CompositeFace.Iface, components)
	return nil
}

// RemoveComponent 使用名称删除组件，同名组件均会删除
func (entity *EntityBehavior) RemoveComponent(name string) {
	compNode, ok := entity.getMutableComponentNode(name)
	if !ok {
		return
	}

	entity.mutableComponentList.TraversalAt(func(node *generic.Node[Component]) bool {
		comp := node.V

		if comp.GetName() != name {
			return false
		}

		if comp.GetFixed() {
			return true
		}

		if comp.GetState() > ComponentState_Alive {
			return true
		}
		comp.setState(ComponentState_Detach)

		_EmitEventComponentMgrRemoveComponent(entity, entity.opts.CompositeFace.Iface, comp)

		node.Escape()
		return true
	}, compNode)
}

// RemoveComponentById 使用组件Id删除组件
func (entity *EntityBehavior) RemoveComponentById(id uid.Id) {
	compNode, ok := entity.getMutableComponentNodeById(id)
	if !ok {
		return
	}

	comp := compNode.V

	if comp.GetFixed() {
		return
	}

	if comp.GetState() > ComponentState_Alive {
		return
	}
	comp.setState(ComponentState_Detach)

	_EmitEventComponentMgrRemoveComponent(entity, entity.opts.CompositeFace.Iface, comp)

	compNode.Escape()
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

func (entity *EntityBehavior) addMutableComponent(name string, comp Component) {
	comp.setFixed(false)
	comp.init(name, entity.opts.CompositeFace.Iface, comp)

	if at, ok := entity.getMutableComponentNode(name); ok {
		entity.mutableComponentList.TraversalAt(func(node *generic.Node[Component]) bool {
			if node.V.GetName() == name {
				at = node
				return true
			}
			return false
		}, at)

		entity.mutableComponentList.InsertAfter(comp, at)

	} else {
		entity.mutableComponentList.PushBack(comp)
	}

	comp.setState(ComponentState_Attach)
}

func (entity *EntityBehavior) getMutableComponentNode(name string) (*generic.Node[Component], bool) {
	var compNode *generic.Node[Component]

	entity.mutableComponentList.Traversal(func(node *generic.Node[Component]) bool {
		if node.V.GetName() == name {
			compNode = node
			return false
		}
		return true
	})

	return compNode, compNode != nil
}

func (entity *EntityBehavior) getMutableComponentNodeById(id uid.Id) (*generic.Node[Component], bool) {
	var compNode *generic.Node[Component]

	entity.mutableComponentList.Traversal(func(node *generic.Node[Component]) bool {
		if node.V.GetId() == id {
			compNode = node
			return false
		}
		return true
	})

	return compNode, compNode != nil
}

func (entity *EntityBehavior) addFixedComponent(name string, comp Component) {
	comp.setFixed(true)
	comp.init(name, entity.opts.CompositeFace.Iface, comp)

	if entity.fixedComponentList.TryAdd(name, comp) {
		comp.setState(ComponentState_Attach)
	}
}

func (entity *EntityBehavior) getFixedComponentNode(name string) (Component, bool) {
	return entity.fixedComponentList.Get(name)
}

func (entity *EntityBehavior) getFixedComponentNodeById(id uid.Id) (Component, bool) {
	var comp Component

	entity.fixedComponentList.Range(func(_ string, v Component) bool {
		if v.GetId() == id {
			comp = v
			return false
		}
		return true
	})

	return comp, comp != nil
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
