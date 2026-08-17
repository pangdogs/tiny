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

package tiny

import (
	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/tiny/runtime"
)

type iWorker interface {
	// Run 启动工作循环并返回终止信号。
	Run() async.Signal
	// Terminate 请求停止并返回终止信号。
	Terminate() async.Signal
	// Terminated 返回终止信号。
	Terminated() async.Signal
	// AdvanceFrames 在 Manual 模式下推进指定帧数。
	AdvanceFrames(frames int64) async.Future
	// AdvanceToFrame 在 Manual 模式下推进到目标帧。
	AdvanceToFrame(frame int64) async.Future
	// AdvanceWhile 在 Manual 模式下持续推进，直到 predicate 返回 false。
	AdvanceWhile(predicate generic.Func1[runtime.Context, bool]) async.Future
}
