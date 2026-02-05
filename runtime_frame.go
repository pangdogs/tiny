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

	"git.golaxy.org/tiny/runtime"
)

type _Frame struct {
	mode                 runtime.FrameMode
	targetFPS            float64
	totalFrames          int64
	curFPS               float64
	curFrames            int64
	runningBeginTime     time.Time
	runningElapseTime    time.Duration
	loopBeginTime        time.Time
	lastLoopElapseTime   time.Duration
	updateBeginTime      time.Time
	lastUpdateElapseTime time.Duration
	statFPSBeginTime     time.Time
	statFPSFrames        int64
	fixedLoopElapseTime  time.Duration
}

// Mode 获取帧模式
func (frame *_Frame) Mode() runtime.FrameMode {
	return frame.mode
}

// TargetFPS 获取目标FPS
func (frame *_Frame) TargetFPS() float64 {
	return frame.targetFPS
}

// CurFPS 获取当前FPS
func (frame *_Frame) CurFPS() float64 {
	return frame.curFPS
}

// TotalFrames 获取运行帧数上限
func (frame *_Frame) TotalFrames() int64 {
	return frame.totalFrames
}

// CurFrames 获取当前帧数
func (frame *_Frame) CurFrames() int64 {
	return frame.curFrames
}

// RunningBeginTime 获取运行开始时间
func (frame *_Frame) RunningBeginTime() time.Time {
	return frame.runningBeginTime
}

// RunningElapseTime 获取运行持续时间
func (frame *_Frame) RunningElapseTime() time.Duration {
	return frame.runningElapseTime
}

// LoopBeginTime 获取当前帧循环开始时间（包含异步调用）
func (frame *_Frame) LoopBeginTime() time.Time {
	return frame.loopBeginTime
}

// LastLoopElapseTime 获取上一帧循环耗时（包含异步调用）
func (frame *_Frame) LastLoopElapseTime() time.Duration {
	return frame.lastLoopElapseTime
}

// UpdateBeginTime 获取当前帧更新开始时间
func (frame *_Frame) UpdateBeginTime() time.Time {
	return frame.updateBeginTime
}

// LastUpdateElapseTime 获取上一次帧更新耗时
func (frame *_Frame) LastUpdateElapseTime() time.Duration {
	return frame.lastUpdateElapseTime
}

func (frame *_Frame) init(mode runtime.FrameMode, targetFPS float64, totalFrames int64) {
	frame.mode = mode
	frame.targetFPS = targetFPS
	frame.totalFrames = totalFrames
}

func (frame *_Frame) setCurFrames(v int64) {
	frame.curFrames = v
}

func (frame *_Frame) runningBegin() {
	now := time.Now()

	frame.curFPS = 0
	frame.curFrames = 0

	frame.statFPSBeginTime = now
	frame.statFPSFrames = 0

	frame.runningBeginTime = now
	frame.runningElapseTime = 0

	frame.loopBeginTime = now
	frame.lastLoopElapseTime = 0

	frame.updateBeginTime = now
	frame.lastUpdateElapseTime = 0

	switch frame.mode {
	case runtime.FrameMode_Simulate, runtime.FrameMode_Manual:
		frame.fixedLoopElapseTime = time.Second / time.Duration(frame.targetFPS)
	}
}

func (frame *_Frame) runningEnd() {
}

func (frame *_Frame) loopBegin() {
	switch frame.mode {
	case runtime.FrameMode_Simulate, runtime.FrameMode_Manual:
		frame.loopBeginTime = frame.runningBeginTime.Add(frame.fixedLoopElapseTime * time.Duration(frame.curFrames))

	default:
		frame.loopBeginTime = time.Now()

		statInterval := frame.loopBeginTime.Sub(frame.statFPSBeginTime).Seconds()
		if statInterval >= 1 {
			frame.curFPS = float64(frame.statFPSFrames) / statInterval
			frame.statFPSBeginTime = frame.loopBeginTime
			frame.statFPSFrames = 0
		}
	}
}

func (frame *_Frame) loopEnd() {
	switch frame.mode {
	case runtime.FrameMode_Simulate, runtime.FrameMode_Manual:
		frame.lastLoopElapseTime = frame.fixedLoopElapseTime
		frame.runningElapseTime += frame.fixedLoopElapseTime

	default:
		frame.lastLoopElapseTime = time.Now().Sub(frame.loopBeginTime)
		frame.runningElapseTime += frame.lastLoopElapseTime
		frame.statFPSFrames++
	}
}

func (frame *_Frame) updateBegin() {
	switch frame.mode {
	case runtime.FrameMode_Simulate, runtime.FrameMode_Manual:
		frame.updateBeginTime = frame.runningBeginTime.Add(frame.fixedLoopElapseTime * time.Duration(frame.curFrames))
	default:
		frame.updateBeginTime = time.Now()
	}
}

func (frame *_Frame) updateEnd() {
	switch frame.mode {
	case runtime.FrameMode_Simulate, runtime.FrameMode_Manual:
		frame.lastUpdateElapseTime = frame.fixedLoopElapseTime
	default:
		frame.lastUpdateElapseTime = time.Now().Sub(frame.updateBeginTime)
	}
}
