package tiny

import (
	"git.golaxy.org/tiny/runtime"
)

func (rt *RuntimeBehavior) loopingManual() {
	frame := runtime.UnsafeFrame(rt.opts.Frame)

	totalFrames := frame.GetTotalFrames()
	gcFrames := int64(rt.opts.GCInterval.Seconds() * float64(frame.GetTargetFPS()))

	var pauseFrame int64

loop:
	for rt.frameLoopBegin(); ; {
		curFrames := frame.GetCurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			break loop
		}

		if curFrames%gcFrames == 0 {
			rt.runGC()
		}

		if frame.GetCurFrames() >= pauseFrame {
		waiting:
			for {
				select {
				case ctrl := <-rt.ctrlChan:
					if ctrl.at {
						if pauseFrame < ctrl.frames {
							pauseFrame = ctrl.frames
						}
					} else {
						pauseFrame += ctrl.frames
					}
					break waiting

				case task, ok := <-rt.processQueue:
					if !ok {
						break loop
					}
					rt.runTask(task)

				case <-rt.ctx.Done():
					break loop
				}
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
