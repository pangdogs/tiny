package ec

import (
	"git.golaxy.org/tiny/event"
	"git.golaxy.org/tiny/utils/uid"
	"reflect"
)

// Deprecated: UnsafeComponent 访问组件内部函数
func UnsafeComponent(comp Component) _UnsafeComponent {
	return _UnsafeComponent{
		Component: comp,
	}
}

type _UnsafeComponent struct {
	Component
}

// Init 初始化
func (uc _UnsafeComponent) Init(name string, entity Entity, composite Component) {
	uc.init(name, entity, composite)
}

// SetId 设置Id
func (uc _UnsafeComponent) SetId(id uid.Id) {
	uc.setId(id)
}

// SetState 设置状态
func (uc _UnsafeComponent) SetState(state ComponentState) {
	uc.setState(state)
}

// SetReflected 设置反射值
func (uc _UnsafeComponent) SetReflected(v reflect.Value) {
	uc.setReflected(v)
}

// GetComposite 获取扩展者
func (uc _UnsafeComponent) GetComposite() Component {
	return uc.getComposite()
}

// GetFixed 获取是否固定
func (uc _UnsafeComponent) GetFixed() bool {
	return uc.getFixed()
}

// SetFixed 设置是否固定
func (uc _UnsafeComponent) SetFixed(fixed bool) {
	uc.setFixed(fixed)
}

// EventComponentDestroySelf 事件：组件销毁自身
func (uc _UnsafeComponent) EventComponentDestroySelf() event.IEvent {
	return uc.eventComponentDestroySelf()
}

// CleanManagedHooks 清理所有的托管hook
func (uc _UnsafeComponent) CleanManagedHooks() {
	uc.cleanManagedHooks()
}
