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

package tiny_test

import (
	"context"
	"log"
	"strings"
	"testing"
	"time"

	"git.golaxy.org/core/define"
	"git.golaxy.org/tiny"
	"git.golaxy.org/tiny/ec"
	"git.golaxy.org/tiny/ec/pt"
	"git.golaxy.org/tiny/runtime"
	"git.golaxy.org/tiny/utils/assertion"
	"git.golaxy.org/tiny/utils/id"
	"github.com/elliotchance/pie/v2"
)

type EntityTest1 struct {
	ec.EntityBehavior
}

func (e *EntityTest1) Awake() {
	log.Printf("EntityTest1 %s Awake", e.Id())
}

func (e *EntityTest1) Start() {
	log.Printf("EntityTest1 %s Start", e.Id())
}

func (e *EntityTest1) Shut() {
	log.Printf("EntityTest1 %s Shut", e.Id())
}

func (e *EntityTest1) Dispose() {
	log.Printf("EntityTest1 %s Dispose", e.Id())
}

type EntityTest2 struct {
	ec.EntityBehavior
}

func (e *EntityTest2) Awake() {
	log.Printf("EntityTest2 %s Awake", e.Id())
}

func (e *EntityTest2) Start() {
	log.Printf("EntityTest2 %s Start", e.Id())
}

func (e *EntityTest2) Shut() {
	log.Printf("EntityTest2 %s Shut", e.Id())
}

func (e *EntityTest2) Dispose() {
	log.Printf("EntityTest2 %s Dispose", e.Id())
}

type ComponentTest1 struct {
	ec.ComponentBehavior
}

func (c *ComponentTest1) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTest1) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTest1) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTest1) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTest2 struct {
	ec.ComponentBehavior
}

func (c *ComponentTest2) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTest2) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTest2) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTest2) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTest3 struct {
	ec.ComponentBehavior
}

func (c *ComponentTest3) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTest3) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTest3) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTest3) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

func Test_ServiceRegisterEntityPT(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				ctx.EntityLib().Declare(
					pt.NewEntityDescriptor("Test1").SetInstance(EntityTest1{}),
					ComponentTest1{},
				)
				ctx.EntityLib().Declare(
					pt.EntityDescriptor{
						Prototype: "Test2",
						Instance:  EntityTest2{},
					},
					ComponentTest1{},
					ComponentTest2{},
				)
				ctx.EntityLib().Declare(
					"Test3",
					ComponentTest1{},
					ComponentTest2{},
					ComponentTest3{},
				)
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

func Test_CreateEntity(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					SetInstance(EntityTest1{}).
					AddComponent(ComponentTest1{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test2").
					SetInstance(EntityTest2{}).
					AddComponent(ComponentTest1{}).
					AddComponent(ComponentTest2{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test3").
					AddComponent(ComponentTest1{}).
					AddComponent(ComponentTest2{}).
					AddComponent(ComponentTest3{}).
					Declare()
			case runtime.RunningEvent_Started:
				tiny.BuildEntity(ctx, "Test1").New()
				tiny.BuildEntity(ctx, "Test2").New()
				tiny.BuildEntity(ctx, "Test3").New()
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

type ComponentTestEnable1 struct {
	ec.ComponentBehavior
}

func (c *ComponentTestEnable1) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
	c.SetEnabled(false)
}

func (c *ComponentTestEnable1) OnEnable() {
	log.Printf("Component %s.%s Enable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable1) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable1) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable1) OnDisable() {
	log.Printf("Component %s.%s Disable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable1) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTestEnable2 struct {
	ec.ComponentBehavior
}

func (c *ComponentTestEnable2) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable2) OnEnable() {
	log.Printf("Component %s.%s Enable", c.Entity().Id(), c.Name())
	c.SetEnabled(false)
}

func (c *ComponentTestEnable2) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable2) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable2) OnDisable() {
	log.Printf("Component %s.%s Disable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable2) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTestEnable3 struct {
	ec.ComponentBehavior
}

func (c *ComponentTestEnable3) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable3) OnEnable() {
	log.Printf("Component %s.%s Enable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable3) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
	c.SetEnabled(false)
}

func (c *ComponentTestEnable3) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable3) OnDisable() {
	log.Printf("Component %s.%s Disable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable3) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTestEnable4 struct {
	ec.ComponentBehavior
}

func (c *ComponentTestEnable4) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable4) OnEnable() {
	log.Printf("Component %s.%s Enable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable4) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable4) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
	c.SetEnabled(false)
}

func (c *ComponentTestEnable4) OnDisable() {
	log.Printf("Component %s.%s Disable", c.Entity().Id(), c.Name())
}

func (c *ComponentTestEnable4) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

func Test_EntityComponentEnable(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestEnable1{}).
					AddComponent(ComponentTestEnable2{}).
					AddComponent(ComponentTestEnable3{}).
					AddComponent(ComponentTestEnable4{}).
					Declare()
			case runtime.RunningEvent_Started:
				tiny.BuildEntity(ctx, "Test1").New()
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

type ComponentTestDynamic1 struct {
	ec.ComponentBehavior

	test2 *ComponentTest2
	test3 *ComponentTest3
}

func (c *ComponentTestDynamic1) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())
}

func (c *ComponentTestDynamic1) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())

	if err := assertion.Inject(c.Entity(), c); err != nil {
		log.Panicln("Inject error:", err)
	}

	log.Println("Inject:", c.test2, c.test3)
}

func (c *ComponentTestDynamic1) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTestDynamic1) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

type ComponentTestDynamic2 struct {
	ec.ComponentBehavior

	test2 *ComponentTest2
	test3 *ComponentTest3
}

func (c *ComponentTestDynamic2) Awake() {
	log.Printf("Component %s.%s Awake", c.Entity().Id(), c.Name())

	if err := assertion.Inject(c.Entity(), c); err != nil {
		log.Panicln("Inject error:", err)
	}

	log.Println("Inject:", c.test2, c.test3)
}

func (c *ComponentTestDynamic2) Start() {
	log.Printf("Component %s.%s Start", c.Entity().Id(), c.Name())
}

func (c *ComponentTestDynamic2) Shut() {
	log.Printf("Component %s.%s Shut", c.Entity().Id(), c.Name())
}

func (c *ComponentTestDynamic2) Dispose() {
	log.Printf("Component %s.%s Dispose", c.Entity().Id(), c.Name())
}

func Test_EntityDynamicComponent(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				ctx.EntityLib().ComponentLib().Declare(ComponentTest2{})
				ctx.EntityLib().ComponentLib().Declare(ComponentTest3{})

				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestDynamic1{}).
					AddComponent(ComponentTest2{}).
					Declare()

				tiny.BuildEntityPT(ctx, "Test2").
					AddComponent(ComponentTestDynamic2{}).
					Declare()

			case runtime.RunningEvent_Started:
				tiny.BuildEntity(ctx, "Test1").New()
				tiny.BuildEntity(ctx, "Test2").New()
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

type ComponentTestParent struct {
	ec.ComponentBehavior
}

func (c *ComponentTestParent) Awake() {
	ec.BindEventTreeNodeAddChild(c.Entity(), c)
	ec.BindEventTreeNodeRemoveChild(c.Entity(), c)
}

func (c *ComponentTestParent) OnTreeNodeAddChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeAddChild %s <- %s", entity.Id(), childId)
}

func (c *ComponentTestParent) OnTreeNodeRemoveChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeRemoveChild %s x- %s", entity.Id(), childId)
}

type ComponentTestChild struct {
	ec.ComponentBehavior
}

func (c *ComponentTestChild) Awake() {
	ec.BindEventTreeNodeAttachParent(c.Entity(), c)
	ec.BindEventTreeNodeDetachParent(c.Entity(), c)
}

func (c *ComponentTestChild) OnTreeNodeAttachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeAttachParent %s -> %s", entity.Id(), parentId)
}

func (c *ComponentTestChild) OnTreeNodeDetachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeDetachParent %s -x %s", entity.Id(), parentId)
}

func PrintEntityTreeForest(entityTree runtime.EntityTree) {
	entityTree.EachChildren(runtime.ForestNodeId, func(entity ec.Entity) {
		PrintEntityTree(entity)
	})
}

func PrintEntityTree(entity ec.Entity, depth ...int) {
	entityTree := runtime.Current(entity).EntityTree()
	if b, _ := entityTree.IsFreedom(entity.Id()); b {
		return
	}

	root := ""

	isRoot, _ := entityTree.IsRoot(entity.Id())
	if isRoot {
		root = "R"
	}

	leaf := ""

	isLeaf, _ := entityTree.IsLeaf(entity.Id())
	if isLeaf {
		leaf = "L"
	}

	_depth := pie.First(depth)

	if isLeaf {
		log.Printf("%s- [%s] %s%s", strings.Repeat(" ", _depth), entity.Id(), root, leaf)
	} else {
		log.Printf("%s+ [%s] %s%s", strings.Repeat(" ", _depth), entity.Id(), root, leaf)
	}

	entityTree.EachChildren(entity.Id(), func(entity ec.Entity) {
		PrintEntityTree(entity, _depth+1)
	})
}

func Test_EntityTree(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChild{}).
					Declare()
			case runtime.RunningEvent_Starting:
				runtime.BindEventEntityTreeAddNode(ctx.EntityTree(), runtime.HandleEventEntityTreeAddNode(func(entityTree runtime.EntityTree, parentId, childId id.Id) {
					var children []id.Id

					entityTree.EachChildren(parentId, func(entity ec.Entity) {
						children = append(children, entity.Id())
					})

					log.Printf("OnEntityTreeAddNode %s: %v + %s", parentId, children, childId)
				}))
				runtime.BindEventEntityTreeRemoveNode(ctx.EntityTree(), runtime.HandleEventEntityTreeRemoveNode(func(entityTree runtime.EntityTree, parentId, childId id.Id) {
					var children []id.Id

					entityTree.EachChildren(parentId, func(entity ec.Entity) {
						children = append(children, entity.Id())
					})

					log.Printf("OnEntityTreeRemoveNode %s: %v - %s", parentId, children, childId)
				}))
				runtime.BindEventEntityTreeMoveNode(ctx.EntityTree(), runtime.HandleEventEntityTreeMoveNode(func(entityTree runtime.EntityTree, childId, fromParentId, toParentId id.Id) {
					log.Printf("OnEntityTreeMoveNode %s: %s => %s", childId, fromParentId, toParentId)
				}))
			case runtime.RunningEvent_Started:
				root, err := tiny.BuildEntity(ctx, "Test1").New()
				if err != nil {
					log.Panicln("new root error:", err)
				}

				err = ctx.EntityTree().MakeRoot(root.Id())
				if err != nil {
					log.Panicln("make root error:", err)
				}

				child1, err := tiny.BuildEntity(ctx, "Test1").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child1 error:", err)
				}

				child2, err := tiny.BuildEntity(ctx, "Test1").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child2 error:", err)
				}

				child3, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child1.Id()).New()
				if err != nil {
					log.Panicln("new child3 error:", err)
				}

				child4, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child3.Id()).New()
				if err != nil {
					log.Panicln("new child4 error:", err)
				}

				child5, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child3.Id()).New()
				if err != nil {
					log.Panicln("new child5 error:", err)
				}

				child6, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child3.Id()).New()
				if err != nil {
					log.Panicln("new child6 error:", err)
				}

				child7, err := tiny.BuildEntity(ctx, "Test1").SetParentId(runtime.ForestNodeId).New()
				if err != nil {
					log.Panicln("new child7 error:", err)
				}

				child8, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child2.Id()).New()
				if err != nil {
					log.Panicln("new child8 error:", err)
				}

				log.Println("1. testing detach node")

				PrintEntityTreeForest(ctx.EntityTree())

				err = ctx.EntityTree().DetachNode(child2.Id())
				if err != nil {
					log.Panicln("detach child2 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("2. testing remove node")

				PrintEntityTreeForest(ctx.EntityTree())

				err = ctx.EntityTree().RemoveNode(child3.Id())
				if err != nil {
					log.Panicln("remove child3 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("3. testing move node")

				PrintEntityTreeForest(ctx.EntityTree())

				err = ctx.EntityTree().MoveNode(child7.Id(), child2.Id())
				if err != nil {
					log.Panicln("move child7 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				_ = child1
				_ = child2
				_ = child3
				_ = child4
				_ = child5
				_ = child6
				_ = child7
				_ = child8
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

type ComponentTestChildDetachInAttaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestChildDetachInAttaching) Awake() {
	ec.BindEventTreeNodeAttachParent(c.Entity(), c)
	ec.BindEventTreeNodeDetachParent(c.Entity(), c)
}

func (c *ComponentTestChildDetachInAttaching) OnTreeNodeAttachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeAttachParent %s -> %s", entity.Id(), parentId)

	err := runtime.Current(entity).EntityTree().DetachNode(entity.Id())
	if err != nil {
		log.Printf("OnTreeNodeAttachParent %s DetachNode failed, %s", entity.Id(), err)
	}
}

func (c *ComponentTestChildDetachInAttaching) OnTreeNodeDetachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeDetachParent %s -x %s", entity.Id(), parentId)
}

type ComponentTestChildRemoveInAttaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestChildRemoveInAttaching) Awake() {
	ec.BindEventTreeNodeAttachParent(c.Entity(), c)
	ec.BindEventTreeNodeDetachParent(c.Entity(), c)
}

func (c *ComponentTestChildRemoveInAttaching) OnTreeNodeAttachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeAttachParent %s -> %s", entity.Id(), parentId)

	err := runtime.Current(entity).EntityTree().RemoveNode(entity.Id())
	if err != nil {
		log.Printf("OnTreeNodeAttachParent %s RemoveNode failed, %s", entity.Id(), err)
	}
}

func (c *ComponentTestChildRemoveInAttaching) OnTreeNodeDetachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeDetachParent %s -x %s", entity.Id(), parentId)
}

type ComponentTestChildDestroyInAttaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestChildDestroyInAttaching) Awake() {
	ec.BindEventTreeNodeAttachParent(c.Entity(), c)
	ec.BindEventTreeNodeDetachParent(c.Entity(), c)
}

func (c *ComponentTestChildDestroyInAttaching) OnTreeNodeAttachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeAttachParent %s -> %s", entity.Id(), parentId)
	entity.Destroy()
}

func (c *ComponentTestChildDestroyInAttaching) OnTreeNodeDetachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeDetachParent %s -x %s", entity.Id(), parentId)
}

type ComponentTestChildDestroyInDetaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestChildDestroyInDetaching) Awake() {
	ec.BindEventTreeNodeAttachParent(c.Entity(), c)
	ec.BindEventTreeNodeDetachParent(c.Entity(), c)
}

func (c *ComponentTestChildDestroyInDetaching) OnTreeNodeAttachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeAttachParent %s -> %s", entity.Id(), parentId)
}

func (c *ComponentTestChildDestroyInDetaching) OnTreeNodeDetachParent(entity ec.Entity, parentId id.Id) {
	log.Printf("OnTreeNodeDetachParent %s -x %s", entity.Id(), parentId)
	entity.Destroy()
}

type ComponentTestParentDestroyInAttaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestParentDestroyInAttaching) Awake() {
	ec.BindEventTreeNodeAddChild(c.Entity(), c)
	ec.BindEventTreeNodeRemoveChild(c.Entity(), c)
}

func (c *ComponentTestParentDestroyInAttaching) OnTreeNodeAddChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeAddChild %s <- %s", entity.Id(), childId)
	entity.Destroy()
}

func (c *ComponentTestParentDestroyInAttaching) OnTreeNodeRemoveChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeRemoveChild %s x- %s", entity.Id(), childId)
}

type ComponentTestParentDestroyInDetaching struct {
	ec.ComponentBehavior
}

func (c *ComponentTestParentDestroyInDetaching) Awake() {
	ec.BindEventTreeNodeAddChild(c.Entity(), c)
	ec.BindEventTreeNodeRemoveChild(c.Entity(), c)
}

func (c *ComponentTestParentDestroyInDetaching) OnTreeNodeAddChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeAddChild %s <- %s", entity.Id(), childId)
}

func (c *ComponentTestParentDestroyInDetaching) OnTreeNodeRemoveChild(entity ec.Entity, childId id.Id) {
	log.Printf("OnTreeNodeRemoveChild %s x- %s", entity.Id(), childId)
	entity.Destroy()
}

func Test_EntityTreeSequence(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChild{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test2").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChildDetachInAttaching{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test3").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChildRemoveInAttaching{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test4").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChildDestroyInAttaching{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test5").
					AddComponent(ComponentTestParent{}).
					AddComponent(ComponentTestChildDestroyInDetaching{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test6").
					AddComponent(ComponentTestParentDestroyInAttaching{}).
					AddComponent(ComponentTestChild{}).
					Declare()
				tiny.BuildEntityPT(ctx, "Test7").
					AddComponent(ComponentTestParentDestroyInDetaching{}).
					AddComponent(ComponentTestChild{}).
					Declare()
			case runtime.RunningEvent_Starting:
				runtime.BindEventEntityTreeAddNode(ctx.EntityTree(), runtime.HandleEventEntityTreeAddNode(func(entityTree runtime.EntityTree, parentId, childId id.Id) {
					var children []id.Id

					entityTree.EachChildren(parentId, func(entity ec.Entity) {
						children = append(children, entity.Id())
					})

					log.Printf("OnEntityTreeAddNode %s: %v + %s", parentId, children, childId)
				}))
				runtime.BindEventEntityTreeRemoveNode(ctx.EntityTree(), runtime.HandleEventEntityTreeRemoveNode(func(entityTree runtime.EntityTree, parentId, childId id.Id) {
					var children []id.Id

					entityTree.EachChildren(parentId, func(entity ec.Entity) {
						children = append(children, entity.Id())
					})

					log.Printf("OnEntityTreeRemoveNode %s: %v - %s", parentId, children, childId)
				}))
				runtime.BindEventEntityTreeMoveNode(ctx.EntityTree(), runtime.HandleEventEntityTreeMoveNode(func(entityTree runtime.EntityTree, childId, fromParentId, toParentId id.Id) {
					log.Printf("OnEntityTreeMoveNode %s: %s => %s", childId, fromParentId, toParentId)
				}))
			case runtime.RunningEvent_Started:
				root, err := tiny.BuildEntity(ctx, "Test1").New()
				if err != nil {
					log.Panicln("new root error:", err)
				}

				err = ctx.EntityTree().MakeRoot(root.Id())
				if err != nil {
					log.Panicln("make root error:", err)
				}

				log.Println("1. testing child detach in attaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child1, err := tiny.BuildEntity(ctx, "Test2").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child1 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("2. testing child remove in attaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child2, err := tiny.BuildEntity(ctx, "Test3").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child2 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("3. testing child destroy in attaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child3, err := tiny.BuildEntity(ctx, "Test4").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child3 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("4. testing child destroy in detaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child4, err := tiny.BuildEntity(ctx, "Test5").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child4 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				ctx.EntityTree().DetachNode(child4.Id())
				log.Printf("%s: state=%s, tree_node_state=%s", child4.Id(), child4.State(), child4.TreeNodeState())

				PrintEntityTreeForest(ctx.EntityTree())

				log.Println("4. testing parent destroy in attaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child5, err := tiny.BuildEntity(ctx, "Test6").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child5 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				child6, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child5.Id()).New()
				if err != nil {
					log.Panicln("new child6 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Printf("%s: state=%s, tree_node_state=%s", child5.Id(), child5.State(), child5.TreeNodeState())
				log.Printf("%s: state=%s, tree_node_state=%s", child6.Id(), child6.State(), child6.TreeNodeState())

				log.Println("5. testing parent destroy in detaching")

				PrintEntityTreeForest(ctx.EntityTree())

				child7, err := tiny.BuildEntity(ctx, "Test7").SetParentId(root.Id()).New()
				if err != nil {
					log.Panicln("new child7 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				child8, err := tiny.BuildEntity(ctx, "Test1").SetParentId(child7.Id()).New()
				if err != nil {
					log.Panicln("new child8 error:", err)
				}

				PrintEntityTreeForest(ctx.EntityTree())

				log.Printf("%s: state=%s, tree_node_state=%s", child7.Id(), child7.State(), child7.TreeNodeState())
				log.Printf("%s: state=%s, tree_node_state=%s", child8.Id(), child8.State(), child8.TreeNodeState())

				ctx.EntityTree().DetachNode(child8.Id())

				PrintEntityTreeForest(ctx.EntityTree())

				log.Printf("%s: state=%s, tree_node_state=%s", child7.Id(), child7.State(), child7.TreeNodeState())
				log.Printf("%s: state=%s, tree_node_state=%s", child8.Id(), child8.State(), child8.TreeNodeState())

				_ = child1
				_ = child2
				_ = child3
				_ = child4
				_ = child5
				_ = child6
				_ = child7
				_ = child8
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}

type ComponentTestFrameUpdate struct {
	ec.ComponentBehavior
}

func (c *ComponentTestFrameUpdate) Update() {
	frame := runtime.Current(c).Frame()
	log.Printf("Component %s.%s Update, fps: %.2f", c.Entity().Id(), c.Name(), frame.CurFPS())
}

func (c *ComponentTestFrameUpdate) LateUpdate() {
	log.Printf("Component %s.%s LateUpdate", c.Entity().Id(), c.Name())
}

func Test_CreateEntityFrameUpdate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestFrameUpdate{}).
					Declare()
			case runtime.RunningEvent_Started:
				for range 10 {
					tiny.BuildEntity(ctx, "Test1").New()
				}
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx).Run()
}

type ComponentTestStressFrameUpdate struct {
	ec.ComponentBehavior
	count int
}

func (c *ComponentTestStressFrameUpdate) Update() {
	c.count++
}

func (c *ComponentTestStressFrameUpdate) LateUpdate() {
	c.count++
}

func Test_CreateEntityStressFrameUpdate(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 120*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				tiny.BuildEntityPT(ctx, "Test1").
					AddComponent(ComponentTestStressFrameUpdate{}).
					Declare()
			case runtime.RunningEvent_FrameLoopBegin:
				for range 200 {
					tiny.BuildEntity(ctx, "Test1").New()
				}
			case runtime.RunningEvent_RunGCBegin:
				log.Printf("fps: %.2f, running_elapse_time: %.3f, last_loop_elapse_time: %.3f, entities: %d",
					ctx.Frame().CurFPS(),
					ctx.Frame().RunningElapseTime().Seconds(),
					ctx.Frame().LastLoopElapseTime().Seconds(),
					ctx.EntityManager().CountEntities())
			}
		}),
	)

	<-tiny.NewRuntime(rtCtx).Run()
}

type RuntimeAddIn1 struct{}

func (RuntimeAddIn1) Init(ctx runtime.Context) {
	log.Println("RuntimeAddIn1 Init")
}

func (RuntimeAddIn1) Shut(ctx runtime.Context) {
	log.Println("RuntimeAddIn1 Shut")
}

func (RuntimeAddIn1) OnContextRunningEvent(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
	log.Println("RuntimeAddIn1 OnContextRunningEvent:", runningEvent)
}

func (RuntimeAddIn1) Hello() {
	log.Println("RuntimeAddIn1 Hello")
}

func NewRuntimeAddIn1(...any) *RuntimeAddIn1 {
	return &RuntimeAddIn1{}
}

var (
	runtimeAddIn1 = define.AddIn(NewRuntimeAddIn1)
)

func Test_RuntimeAddIn(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)

	rtCtx := runtime.NewContext(
		runtime.With.Context(ctx),
		runtime.With.RunningEventCB(func(ctx runtime.Context, runningEvent runtime.RunningEvent, args ...any) {
			switch runningEvent {
			case runtime.RunningEvent_Birth:
				runtimeAddIn1.Install(ctx)
			case runtime.RunningEvent_Started:
				runtimeAddIn1.Require(ctx).Hello()
			}
			log.Println("runtime event:", runningEvent, args)
		}),
	)

	<-tiny.NewRuntime(rtCtx, tiny.With.Runtime.Frame(tiny.With.Frame.Enabled(false))).Run()
}
