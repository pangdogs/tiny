package tiny

import (
	"time"
)

func (rt *RuntimeBehavior) loopingNoFrame() {
	gcTicker := time.NewTicker(rt.opts.GCInterval)
	defer gcTicker.Stop()

loop:
	for {
		select {
		case task, ok := <-rt.processQueue:
			if !ok {
				break loop
			}
			rt.runTask(task)

		case <-gcTicker.C:
			rt.runGC()

		case <-rt.ctx.Done():
			break loop
		}
	}

	close(rt.processQueue)

	for {
		select {
		case task, ok := <-rt.processQueue:
			if !ok {
				return
			}
			rt.runTask(task)

		default:
			return
		}
	}
}
