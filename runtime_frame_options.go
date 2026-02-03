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

	"git.golaxy.org/core/utils/exception"
	"git.golaxy.org/core/utils/option"
	"git.golaxy.org/tiny/runtime"
)

// FrameOptions 帧的所有选项
type FrameOptions struct {
	Enable      bool              // 是否启用帧
	Mode        runtime.FrameMode // 帧模式
	TargetFPS   float64           // 目标FPS
	TotalFrames int64             // 运行帧数上限
}

type _FrameOption struct{}

// Default 默认值
func (_FrameOption) Default() option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		With.Frame.Enable(true).Apply(options)
		With.Frame.Mode(runtime.FrameMode_RealTime).Apply(options)
		With.Frame.TargetFPS(30).Apply(options)
		With.Frame.TotalFrames(0).Apply(options)
	}
}

// Enable 是否启用帧
func (_FrameOption) Enable(b bool) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		options.Enable = b
	}
}

// Mode 帧模式
func (_FrameOption) Mode(m runtime.FrameMode) option.Setting[FrameOptions] {
	return func(o *FrameOptions) {
		o.Mode = m
	}
}

// TargetFPS 目标FPS
func (_FrameOption) TargetFPS(fps float64) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		if fps <= 0 {
			exception.Panicf("%w: %w: TargetFPS less equal 0 is invalid", runtime.ErrFrame, exception.ErrArgs)
		}
		options.TargetFPS = math.Round(fps)
	}
}

// TotalFrames 运行帧数上限
func (_FrameOption) TotalFrames(v int64) option.Setting[FrameOptions] {
	return func(options *FrameOptions) {
		if v < 0 {
			exception.Panicf("%w: %w: TotalFrames less 0 is invalid", runtime.ErrFrame, exception.ErrArgs)
		}
		options.TotalFrames = v
	}
}
