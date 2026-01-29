package tiny

import (
	"git.golaxy.org/tiny/runtime"
)

func (rt *RuntimeBehavior) loopingManual() {
	frame := runtime.UnsafeFrame(rt.options.Frame)

	totalFrames := frame.GetTotalFrames()
	gcFrames := int64(rt.options.GCInterval.Seconds() * frame.GetTargetFPS())

	var curCtrl _Ctrl

loop:
	for rt.frameLoopBegin(); ; {
		curFrames := frame.GetCurFrames()

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

				case task, ok := <-rt.taskQueue:
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
	close(rt.taskQueue)

loopEnding:
	for {
		select {
		case task, ok := <-rt.taskQueue:
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
