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

// Package assertion 提供基于 Tiny 实体组件模型的组件组合与注入辅助。
/*
Package assertion 通过反射从 ec.Entity 中提取组件，并按字段类型、组件名或原型名
注入到结构体中。

它适合把多个组件组合成临时视图，或者在组件启动时补齐依赖字段。字段可通过
`ec:"name,prototype"` tag 指定组件名或组件原型；如果目标原型已在当前 Runtime 的
组件库中注册，但实体上尚未存在对应组件，Inject 还可以按原型动态创建组件。

未找到匹配组件时，对应字段保持零值，不视为错误。查询可能触发组件首次访问 Awake，
动态创建也会修改实体，因此这些操作必须在实体所属 Runtime 的运行 goroutine 中执行。
该包依赖反射并可能产生分配，适合装配期、启动期或测试，不应在帧更新热点中反复调用。
*/
package assertion
