package utils

import "github.com/go-gl/glfw/v3.2/glfw"

var (
	t0Value       float64       // Set the initial time to now
	fpsFrameCount int = 0          // Set the initial FPS frame count to 0
	fps           float64 = 0.0 // Set the initial FPS value to 0.0
)
func InitFPS() {
	t0Value = glfw.GetTime()
}

func CalcFPS(theTimeInterval float64) float64 {
	// Get the current time in seconds since the program started (non-static, so executed every time)
	currentTime := glfw.GetTime()

	// Calculate and display the FPS every specified time interval
	if (currentTime - t0Value) > 1.0 {
		// Calculate the FPS as the number of frames divided by the interval in seconds
		fps = float64(fpsFrameCount) / (currentTime - t0Value)

		// Reset the FPS frame counter and set the initial time to be now
		fpsFrameCount = 0
		t0Value = glfw.GetTime()
	} else {
		// FPS calculation time interval hasn't elapsed yet? Simply increment the FPS frame counter
		fpsFrameCount++
	}

	// Return the current FPS - doesn't have to be used if you don't want it!
	return fps
}