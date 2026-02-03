package tiny

func (rt *RuntimeBehavior) loopingSimulate() {
	frame := rt.frame

	totalFrames := frame.GetTotalFrames()
	gcFrames := int64(rt.options.GCInterval.Seconds() * frame.GetTargetFPS())

loop:
	for rt.frameLoopBegin(); ; {
		curFrames := frame.GetCurFrames()

		if totalFrames > 0 && curFrames >= totalFrames {
			break
		}

		select {
		case <-rt.ctx.Done():
			break loop
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
