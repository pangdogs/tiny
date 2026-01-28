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
	"time"

	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/tiny/runtime"
)

type iWorker interface {
	// Run 运行
	Run() async.AsyncRet
	// Play 播放指定时长
	Play(delta time.Duration) error
	// PlayAt 播放至指定位置
	PlayAt(at time.Duration) error
	// PlayFrames 播放指定帧数
	PlayFrames(delta int64) error
	// PlayAtFrames 播放至指定帧数
	PlayAtFrames(at int64) error
	// PlayAtFunc 播放至函数指定位置
	PlayAtFunc(fun generic.Func1[runtime.Context, bool]) error
	// Terminate 停止
	Terminate() async.AsyncRet
	// Terminated 已停止
	Terminated() async.AsyncRet
}
