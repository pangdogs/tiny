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
	"fmt"

	"git.golaxy.org/core/utils/async"
	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/tiny/runtime"
)

var (
	ErrManualFrameMode   = fmt.Errorf("%w: manual frame mode required", ErrRuntime)
	ErrRuntimeNotRunning = fmt.Errorf("%w: not running", ErrRuntime)
)

// AdvanceFrames 在 Manual 模式下推进 frames 帧，并在完成后返回当前帧号。
func (rt *RuntimeBehavior) AdvanceFrames(frames int64) async.Future {
	if frames <= 0 {
		return async.Rejected(fmt.Errorf("%w: %w: frames must be greater than 0", ErrRuntime, ErrArgs))
	}
	return rt.enqueueManualAdvance(func(ctx runtime.Context) int64 {
		for range frames {
			if !rt.advanceFrame(ctx) {
				break
			}
		}
		return rt.frame.CurFrames()
	})
}

// AdvanceToFrame 在 Manual 模式下推进到 target 帧，并在完成后返回当前帧号。
// target 不大于当前帧时不执行更新。
func (rt *RuntimeBehavior) AdvanceToFrame(target int64) async.Future {
	if target < 0 {
		return async.Rejected(fmt.Errorf("%w: %w: target must be greater than or equal to 0", ErrRuntime, ErrArgs))
	}
	return rt.enqueueManualAdvance(func(ctx runtime.Context) int64 {
		for rt.frame.CurFrames() < target {
			if !rt.advanceFrame(ctx) {
				break
			}
		}
		return rt.frame.CurFrames()
	})
}

// AdvanceWhile 在 Manual 模式下逐帧调用 predicate；返回 false 时停止。
// Future 完成值为停止后的当前帧号。
func (rt *RuntimeBehavior) AdvanceWhile(predicate generic.Func1[runtime.Context, bool]) async.Future {
	if predicate == nil {
		return async.Rejected(fmt.Errorf("%w: %w: predicate is nil", ErrRuntime, ErrArgs))
	}
	return rt.enqueueManualAdvance(func(ctx runtime.Context) int64 {
		for predicate.UnsafeCall(ctx) {
			if !rt.advanceFrame(ctx) {
				break
			}
		}
		return rt.frame.CurFrames()
	})
}

func (rt *RuntimeBehavior) enqueueManualAdvance(advance func(runtime.Context) int64) async.Future {
	if rt.options.Frame.Mode != FrameMode_Manual || rt.frame == nil {
		return async.Rejected(ErrManualFrameMode)
	}
	if !rt.isRunning.Load() {
		return async.Rejected(ErrRuntimeNotRunning)
	}
	if err := rt.ctx.Err(); err != nil {
		return async.Rejected(err)
	}

	return rt.taskQueue.enqueueManualFrame(rt.ctx.ExecutorID(), func(ctx runtime.Context, _ ...any) async.Result {
		return async.NewResult(advance(ctx), nil)
	})
}

func (rt *RuntimeBehavior) advanceFrame(ctx runtime.Context) bool {
	if ctx.Err() != nil {
		return false
	}
	if total := rt.frame.TotalFrames(); total > 0 && rt.frame.CurFrames() >= total {
		rt.Terminate()
		return false
	}

	rt.runFrame()

	if total := rt.frame.TotalFrames(); total > 0 && rt.frame.CurFrames() >= total {
		rt.Terminate()
	}
	return true
}

func (rt *RuntimeBehavior) runFrame() {
	rt.frameLoopBegin()
	rt.frameLoopEnd()
}
