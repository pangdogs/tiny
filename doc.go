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

// Package tiny 提供面向实时计算场景的轻量 Actor + EC 运行内核。
/*
Package tiny 把 runtime、ec 以及 Core 的 event、extension、async 等基础设施
组合为独立的进程内运行模型。

典型流程如下：

  1. 用 runtime.NewContext 创建独立运行时上下文。
  2. 用 BuildEntityPT 声明实体原型。
  3. 用 NewRuntime 绑定并启动运行时。
  4. 用 BuildEntity 创建实体，由 Runtime 串行推进组件生命周期。

此外，根包还提供：

  - 实体、组件、插件的生命周期接口；
  - Submit/Post Actor 邮箱调度，以及保留 Delegate/DelegateVoid 的对应变体；
  - Runtime Scope、懒加载的对象级 Scope、Spawn 后台任务与 ContinueOn Runtime 续体；
  - After、Every 等定时和连续流入口；
  - Realtime、Manual 和 Disabled 三种帧推进模式；
  - Runtime、Frame 与 TaskQueue 的选项构造器；
  - 面向框架集成与高级扩展场景的 unsafe 辅助入口。

Tiny 不包含 Service、RPC、网络、配置或服务发现。Entity 与 Component 不会预创建
goroutine、Context、Signal 或 AsyncScope；对象级 AsyncScope 仅在首次访问时创建，
并在相应对象死亡时关闭。
*/
package tiny
