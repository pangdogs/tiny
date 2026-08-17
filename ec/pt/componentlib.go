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

package pt

import (
	"reflect"

	"git.golaxy.org/core/utils/exception"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/core/utils/types"
	"git.golaxy.org/tiny/ec"
)

// ComponentLib 是 Runtime 本地的组件原型注册表，不提供并发保护。
type ComponentLib interface {
	// Declare 声明组件原型；同一具名类型重复声明时返回已有原型。
	Declare(comp any) ec.ComponentPT
	// Get 按完整原型名查询组件原型。
	Get(prototype string) (ec.ComponentPT, bool)
	// List 返回当前全部组件原型的副本。
	List() []ec.ComponentPT

	IComponentLibEventTab
}

// NewComponentLib 创建独立的空组件原型库。
func NewComponentLib() ComponentLib {
	return &_ComponentLib{
		compPTNameIndex: map[string]int{},
	}
}

type _ComponentLib struct {
	compPTNameIndex map[string]int
	compPTList      generic.FreeList[ec.ComponentPT]

	componentLibEventTab
}

// Declare 声明组件原型；同一具名类型重复声明时返回已有原型。
//
// comp 可以是组件值或 reflect.Type。匿名类型、nil 以及未实现 ec.Component 的类型
// 会导致 panic。
func (lib *_ComponentLib) Declare(comp any) ec.ComponentPT {
	if comp == nil {
		exception.Panicf("%w: %w: comp is nil", ErrPt, exception.ErrArgs)
	}

	compRT, ok := comp.(reflect.Type)
	if !ok {
		compRT = reflect.TypeOf(comp)
	}

	for compRT.Kind() == reflect.Pointer {
		compRT = compRT.Elem()
	}

	if compRT.Name() == "" {
		exception.Panicf("%w: anonymous component not allowed", ErrPt)
	}

	prototype := types.FullNameRT(compRT)

	if !reflect.PointerTo(compRT).Implements(reflect.TypeFor[ec.Component]()) {
		exception.Panicf("%w: component %q not implement ec.Component", ErrPt, prototype)
	}

	if compPTIdx, ok := lib.compPTNameIndex[prototype]; ok {
		return lib.compPTList.Get(compPTIdx).V
	}

	compPT := &_Component{
		prototype:  prototype,
		instanceRT: compRT,
	}
	compPT.builtin = &ec.BuiltinComponent{PT: compPT, Offset: -1}

	lib.compPTNameIndex[prototype] = lib.compPTList.PushBack(compPT).Index()

	_EmitEventComponentLibDeclareComponentPT(lib, compPT)

	return compPT
}

// Get 按完整原型名查询组件原型。
func (lib *_ComponentLib) Get(prototype string) (ec.ComponentPT, bool) {
	compPTIdx, ok := lib.compPTNameIndex[prototype]
	if !ok {
		return nil, false
	}
	return lib.compPTList.Get(compPTIdx).V, true
}

// List 返回当前全部组件原型的副本。
func (lib *_ComponentLib) List() []ec.ComponentPT {
	return lib.compPTList.ToSlice()
}
