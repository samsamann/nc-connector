package util

import "runtime"

// GetFuncName returns the name of the function that called this func.
func GetFuncName() string {
	return getFrame(1).Function
}

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2
	programCounters := make([]uintptr, targetFrameIndex+1)

	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for next, frameIndex := true, 0; next; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, next = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
				break
			}
		}
	}

	return frame
}
