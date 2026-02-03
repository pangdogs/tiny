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
	"context"
	"fmt"
	"time"

	"git.golaxy.org/core/utils/generic"
	"git.golaxy.org/tiny/runtime"
)

var (
	ErrCtrlChanClosed = fmt.Errorf("%w: ctrl chan is closed", ErrRuntime) // 运行控制已关闭
)

type _CtrlMode int32

const (
	_CtrlMode_Pending _CtrlMode = iota
	_CtrlMode_FrameDelta
	_CtrlMode_FrameAt
	_CtrlMode_IfContinue
)

type _Ctrl struct {
	mode         _CtrlMode
	frames       int64
	continueFunc generic.Func1[runtime.Context, bool]
}

// Play 播放指定时长
func (rt *RuntimeBehavior) Play(ctx context.Context, delta time.Duration) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if rt.ctrlChan == nil {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	ctrl := _Ctrl{
		mode:   _CtrlMode_FrameDelta,
		frames: int64(delta.Seconds() * rt.frame.GetTargetFPS()),
	}

	select {
	case rt.ctrlChan <- ctrl:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// PlayAt 播放至指定位置
func (rt *RuntimeBehavior) PlayAt(ctx context.Context, at time.Duration) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if rt.ctrlChan == nil {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	ctrl := _Ctrl{
		mode:   _CtrlMode_FrameAt,
		frames: int64(at.Seconds() * rt.frame.GetTargetFPS()),
	}

	select {
	case rt.ctrlChan <- ctrl:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// PlayFrames 播放指定帧数
func (rt *RuntimeBehavior) PlayFrames(ctx context.Context, delta int64) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if rt.ctrlChan == nil {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	ctrl := _Ctrl{
		mode:   _CtrlMode_FrameDelta,
		frames: delta,
	}

	select {
	case rt.ctrlChan <- ctrl:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// PlayFramesAt 播放至指定帧数
func (rt *RuntimeBehavior) PlayFramesAt(ctx context.Context, at int64) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if rt.ctrlChan == nil {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	ctrl := _Ctrl{
		mode:   _CtrlMode_FrameAt,
		frames: at,
	}

	select {
	case rt.ctrlChan <- ctrl:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}

// PlayIfContinue 指定函数判断是否继续播放
func (rt *RuntimeBehavior) PlayIfContinue(ctx context.Context, continueFunc generic.Func1[runtime.Context, bool]) (err error) {
	if ctx == nil {
		ctx = context.Background()
	}

	if rt.ctrlChan == nil {
		return ErrCtrlChanClosed
	}

	defer func() {
		if panicInfo := recover(); panicInfo != nil {
			err = ErrCtrlChanClosed
		}
	}()

	ctrl := _Ctrl{
		mode:         _CtrlMode_IfContinue,
		continueFunc: continueFunc,
	}

	select {
	case rt.ctrlChan <- ctrl:
		break
	case <-ctx.Done():
		return ctx.Err()
	}

	return nil
}
