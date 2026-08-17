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
	"math"
	"time"

	"git.golaxy.org/core/utils/exception"
	"git.golaxy.org/core/utils/option"
	"git.golaxy.org/tiny/runtime"
)

// FrameMode 定义 Runtime 如何驱动帧更新。
type FrameMode int8

const (
	FrameMode_Disabled FrameMode = iota // 不产生 Update/LateUpdate。
	FrameMode_Realtime                  // 按目标 FPS 自动推进。
	FrameMode_Manual                    // 仅由 Advance 系列方法推进。
)

// FrameOptions 定义运行时帧循环的选项。
type FrameOptions struct {
	Mode        FrameMode // 帧推进模式。
	TargetFPS   float64   // 目标 FPS；设置时会四舍五入为整数值。
	TotalFrames int64     // 最大运行帧数；0 表示不限制。
}

type _FrameOption struct{}

// Default 返回帧选项的默认设置。
func (_FrameOption) Default() option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		With.Frame.Mode(FrameMode_Realtime).Apply(options)
		With.Frame.TargetFPS(30).Apply(options)
		With.Frame.TotalFrames(0).Apply(options)
	}
}

// Mode 设置帧推进模式。
func (_FrameOption) Mode(mode FrameMode) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		switch mode {
		case FrameMode_Disabled, FrameMode_Realtime, FrameMode_Manual:
			options.Mode = mode
		default:
			exception.Panicf("%w: %w: invalid frame mode %d", runtime.ErrFrame, exception.ErrArgs, mode)
		}
	}
}

// TargetFPS 设置目标 FPS；fps 会被四舍五入为 time.Ticker 可表示的整数值。
func (_FrameOption) TargetFPS(fps float64) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		fps = math.Round(fps)
		if math.IsNaN(fps) || math.IsInf(fps, 0) || fps < 1 || fps > float64(time.Second) {
			exception.Panicf("%w: %w: TargetFPS must round to a value between 1 and %d", runtime.ErrFrame, exception.ErrArgs, time.Second)
		}
		options.TargetFPS = fps
	}
}

// TotalFrames 设置最大运行帧数；0 表示不限制，负值会导致 panic。
func (_FrameOption) TotalFrames(v int64) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		if v < 0 {
			exception.Panicf("%w: %w: TotalFrames must be greater than or equal to 0", runtime.ErrFrame, exception.ErrArgs)
		}
		options.TotalFrames = v
	}
}
