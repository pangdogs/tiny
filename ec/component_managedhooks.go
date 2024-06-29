package ec

import (
	"git.golaxy.org/tiny/event"
	"golang.org/x/exp/slices"
)

// ManagedHooks 托管hook，在组件销毁时自动解绑定
func (comp *ComponentBehavior) ManagedHooks(hooks ...event.Hook) {
	comp.managedHooks = slices.DeleteFunc(comp.managedHooks, func(hook event.Hook) bool {
		return !hook.IsBound() || slices.Contains(hooks, hook)
	})
	comp.managedHooks = append(comp.managedHooks, hooks...)
}

func (comp *ComponentBehavior) cleanManagedHooks() {
	event.Clean(comp.managedHooks)
	comp.managedHooks = nil
}
