package tiny

import (
	"git.golaxy.org/tiny/runtime"
)

func (rt *RuntimeBehavior) loopingManual() {
	frame := runtime.UnsafeFrame(rt.opts.Frame)

	totalFrames := frame.GetTotalFrames()
	gcFrames := int64(rt.opts.GCInterval.Seconds() * float64(frame.GetTargetFPS()))

	var curCtrl _Ctrl

loop:
	for rt.frameLoopBegin(); ; {
		curFrames := frame.GetCurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			break loop
		}

		if curFrames%gcFrames == 0 {
			rt.runGC()
		}

	pause:
		switch curCtrl.mode {
		case _NoCtrl:
			for {
				select {
				case ctrl := <-rt.ctrlChan:
					curCtrl = ctrl

					if curCtrl.mode == _FrameDelta {
						curCtrl.mode = _FrameAt
						curCtrl.frames += curFrames
					}

					goto pause

				case task, ok := <-rt.processQueue:
					if !ok {
						break loop
					}
					rt.runTask(task)

				case <-rt.ctx.Done():
					break loop
				}
			}
		case _FrameAt:
			if curFrames >= curCtrl.frames {
				curCtrl = _Ctrl{}
				goto pause
			}
		case _FuncAt:
			if curCtrl.fun.Exec() {
				curCtrl = _Ctrl{}
				goto pause
			}
		}

		rt.frameLoop(nil)
	}

	close(rt.ctrlChan)
	close(rt.processQueue)

loopEnding:
	for {
		select {
		case task, ok := <-rt.processQueue:
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
