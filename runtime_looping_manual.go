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

func (rt *RuntimeBehavior) loopingManual() {
	frame := rt.frame

	totalFrames := frame.TotalFrames()
	gcFrames := int64(rt.options.GCInterval.Seconds() * frame.TargetFPS())

	var curCtrl _Ctrl
	taskOut := rt.taskQueue.out()

loop:
	for rt.frameLoopBegin(); ; {
		curFrames := frame.CurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			break loop
		}

		select {
		case <-rt.ctx.Done():
			break loop
		default:
		}

		if curFrames%gcFrames == 0 {
			rt.runGC()
		}

	pause:
		switch curCtrl.mode {
		case _CtrlMode_Pending:
			for {
				select {
				case ctrl := <-rt.ctrlChan:
					if ctrl.mode == _CtrlMode_FrameDelta {
						ctrl.mode = _CtrlMode_FrameAt
						ctrl.frames += curFrames
					}
					curCtrl = ctrl
					goto pause

				case task, ok := <-taskOut:
					if !ok {
						break loop
					}
					rt.runTask(task)

				case <-rt.ctx.Done():
					break loop
				}
			}
		case _CtrlMode_FrameAt:
			if curFrames >= curCtrl.frames {
				curCtrl = _Ctrl{}
				goto pause
			}
		case _CtrlMode_IfContinue:
			if !curCtrl.continueFunc.UnsafeCall(rt.ctx) {
				curCtrl = _Ctrl{}
				goto pause
			}
		}

		rt.frameLoop(nil)
	}

	close(rt.ctrlChan)
	rt.taskQueue.close()

loopEnding:
	for {
		select {
		case task, ok := <-taskOut:
			if !ok {
				break loopEnding
			}
			rt.runTask(task)

		default:
			break loopEnding
		}
	}

	rt.runGC()
	rt.frameLoopEnd()
}
