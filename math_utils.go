package main

import "math"

func angleOfLine(x1, y1, x2, y2 float64) float64 {

	xDiff := x2 - x1
	yDiff := y2 - y1
	return math.Atan2(yDiff, xDiff)
}
