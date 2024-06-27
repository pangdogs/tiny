package tiny

import "git.golaxy.org/tiny/runtime"

func (rt *RuntimeBehavior) loopingSimulate() {
	frame := runtime.UnsafeFrame(rt.opts.Frame)

	totalFrames := frame.GetTotalFrames()
	gcFrames := int64(rt.opts.GCInterval.Seconds() * float64(frame.GetTargetFPS()))

	for rt.frameLoopBegin(); ; {
		curFrames := frame.GetCurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			break
		}

		select {
		case <-rt.ctx.Done():
			break
		default:
		}

		if curFrames%gcFrames == 0 {
			rt.runGC()
		}

		rt.frameLoop(nil)
	}

	rt.runGC()
	rt.frameLoopEnd()
}
