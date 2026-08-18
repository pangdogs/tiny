/*
 * This file is part of Golaxy Distributed Service Development Framework.
 *
 * Golaxy Distributed Service Development Framework is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 2.1 of the License, or
 * (at your option) any later version.
 *
 * Golaxy Distributed Service Development Framework is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with Golaxy Distributed Service Development Framework. If not, see <http://www.gnu.org/licenses/>.
 *
 * Copyright (c) 2024 pangdogs.
 */

package ec

import (
	"context"
	"fmt"

	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/corectx"
	"git.golaxy.org/core/utils/iface"
	"git.golaxy.org/tiny/utils/id"
)

// ConcurrentEntity 暴露可跨协程安全读取的实体信息、上下文和生命周期 Scope。
//
// 实体成功加入 Runtime 后，才能跨 goroutine 使用此接口；提前调用依赖 Runtime
// Context 的方法属于未定义行为。AsyncScope 和 String 在 Context 绑定前返回空值。
//
// 组件管理、实体树和 Destroy 等操作仍须通过所属 Runtime 的运行协程执行。
type ConcurrentEntity interface {
	iConcurrentEntity
	corectx.ConcurrentContextProvider
	corectx.AsyncScopeProvider
	fmt.Stringer

	// ID 返回实体 ID。
	ID() id.ID
	// PT 返回实体原型。
	PT() EntityPT
}

// runtimeContext 汇集 Entity 绑定所属 Runtime 时需要的最小能力。
// 使用结构接口避免 ec 反向依赖 runtime 包。
type runtimeContext interface {
	context.Context
	corectx.CurrentContextProvider
}

type iConcurrentEntity interface {
	getInstance() Entity
	setContext(rtCtx runtimeContext)
}

// entityAsyncScopeState 发布后不可变；nil 指针表示尚未创建且仍可用。
type entityAsyncScopeState struct {
	scope  *async.Scope
	closed bool
}

var closedEntityAsyncScopeState = &entityAsyncScopeState{closed: true}

// ConcurrentContextCache 返回实体所属 Runtime 的并发上下文接口缓存。
func (entity *EntityBehavior) ConcurrentContextCache() iface.Cache {
	return entity.runtimeCtx.ConcurrentContextCache()
}

// AsyncScope 返回绑定实体生命周期的懒加载后台任务作用域。
// Scope 在实体进入 Dead 时关闭；Runtime Context 尚未绑定时返回 nil。
// 实体已关闭后首次访问会返回已关闭的 Scope。
func (entity *EntityBehavior) AsyncScope() *async.Scope {
	for {
		state := entity.asyncScope.Load()
		if state != nil && state.scope != nil {
			return state.scope
		}

		if entity.runtimeCtx == nil {
			return nil
		}

		asyncScope := async.NewScope(entity.runtimeCtx)
		closed := state != nil && state.closed
		if closed {
			asyncScope.Close()
		}

		newState := &entityAsyncScopeState{
			scope:  asyncScope,
			closed: closed,
		}
		if entity.asyncScope.CompareAndSwap(state, newState) {
			return asyncScope
		}

		asyncScope.Close()
	}
}

// String 返回实体 ID 的十进制文本；Runtime Context 尚未绑定时返回空字符串。
func (entity *EntityBehavior) String() string {
	if entity.runtimeCtx == nil {
		return ""
	}
	return entity.ID().String()
}

func (entity *EntityBehavior) getInstance() Entity {
	return entity.options.InstanceFace.Iface
}

func (entity *EntityBehavior) setContext(rtCtx runtimeContext) {
	if entity.runtimeCtx != nil {
		return
	}
	entity.runtimeCtx = rtCtx
}

func (entity *EntityBehavior) closeAsyncScope() {
	for {
		state := entity.asyncScope.Load()
		if state == nil {
			if entity.asyncScope.CompareAndSwap(nil, closedEntityAsyncScopeState) {
				return
			}
			continue
		}
		if state.closed {
			return
		}

		// 先关闭实际 Scope 再发布关闭状态，避免读取方看到关闭标记时 Scope 仍可接收任务。
		state.scope.Close()
		closedState := &entityAsyncScopeState{
			scope:  state.scope,
			closed: true,
		}
		if entity.asyncScope.CompareAndSwap(state, closedState) {
			return
		}
	}
}
